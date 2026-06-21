package lww

import (
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/utils/hlc"
)

const envelopeVersion byte = 1

const flagDeleted byte = 1 << 0

const headerSize = 16

type Envelope struct {
	Timestamp hlc.Timestamp
	Deleted   bool
	Payload   []byte
}

func (e Envelope) Marshal() []byte {
	node := []byte(e.Timestamp.NodeID)

	buf := make([]byte, headerSize, headerSize+len(node)+len(e.Payload))

	buf[0] = envelopeVersion
	if e.Deleted {
		buf[1] |= flagDeleted
	}
	binary.BigEndian.PutUint64(buf[2:10], uint64(e.Timestamp.WallTime))
	binary.BigEndian.PutUint32(buf[10:14], e.Timestamp.Logical)
	binary.BigEndian.PutUint16(buf[14:16], uint16(len(node)))

	buf = append(buf, node...)
	buf = append(buf, e.Payload...)
	return buf
}

func Unmarshal(data []byte) (Envelope, error) {
	if len(data) < headerSize {
		return Envelope{}, fmt.Errorf("lww: envelope too short: %d bytes", len(data))
	}
	if data[0] != envelopeVersion {
		return Envelope{}, fmt.Errorf("lww: unsupported envelope version: %d", data[0])
	}

	nodeLen := int(binary.BigEndian.Uint16(data[14:16]))
	if len(data) < headerSize+nodeLen {
		return Envelope{}, fmt.Errorf("lww: envelope truncated")
	}

	e := Envelope{
		Timestamp: hlc.Timestamp{
			WallTime: int64(binary.BigEndian.Uint64(data[2:10])),
			Logical:  binary.BigEndian.Uint32(data[10:14]),
			NodeID:   string(data[headerSize : headerSize+nodeLen]),
		},
		Deleted: data[1]&flagDeleted != 0,
	}

	if payload := data[headerSize+nodeLen:]; len(payload) > 0 {
		e.Payload = append([]byte(nil), payload...)
	}

	return e, nil
}

func Apply(txn *badger.Txn, key []byte, incoming Envelope) (bool, error) {
	item, err := txn.Get(key)
	switch {
	case errors.Is(err, badger.ErrKeyNotFound):
	case err != nil:
		return false, err
	default:
		local, err := readEnvelope(item)
		if err != nil {
			return false, err
		}
		if !incoming.Timestamp.After(local.Timestamp) {
			return false, nil
		}
	}

	if err := txn.Set(key, incoming.Marshal()); err != nil {
		return false, err
	}
	return true, nil
}

func Read(txn *badger.Txn, key []byte) (Envelope, bool, error) {
	item, err := txn.Get(key)
	if errors.Is(err, badger.ErrKeyNotFound) {
		return Envelope{}, false, nil
	}
	if err != nil {
		return Envelope{}, false, err
	}

	e, err := readEnvelope(item)
	if err != nil {
		return Envelope{}, false, err
	}
	if e.Deleted {
		return e, false, nil
	}
	return e, true, nil
}

func readEnvelope(item *badger.Item) (Envelope, error) {
	var e Envelope
	err := item.Value(func(val []byte) error {
		var uerr error
		e, uerr = Unmarshal(val)
		return uerr
	})
	return e, err
}
