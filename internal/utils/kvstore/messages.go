package kvstore

import (
	"fmt"

	"bufio"
	"encoding/binary"
	"io"

	"google.golang.org/protobuf/proto"

	"github.com/dgraph-io/badger/v4"
	"github.com/dgraph-io/badger/v4/pb"
)

type message struct {
	replyTo chan<- error
	operation
}

type operation interface {
	Handle(db *badger.DB) error
}

type backupOperation struct {
	w     io.Writer
	since uint64
}

type restoreOperation struct {
	r io.Reader
}

type mergeOperation struct {
	r       io.Reader
	mergeFn func(txn *badger.Txn, kv *pb.KV) error
}

type viewOperation struct {
	txnFn func(txn *badger.Txn) error
}

type updateOperation struct {
	txnFn func(txn *badger.Txn) error
}

type latestVersionOperation struct {
	version uint64
}

type dropAllOperation struct{}

var _ operation = (*backupOperation)(nil)
var _ operation = (*restoreOperation)(nil)
var _ operation = (*mergeOperation)(nil)
var _ operation = (*viewOperation)(nil)
var _ operation = (*updateOperation)(nil)
var _ operation = (*latestVersionOperation)(nil)
var _ operation = (*dropAllOperation)(nil)

func (m *backupOperation) Handle(db *badger.DB) error {
	var err error
	stream := db.NewStream()
	stream.NumGo = 1
	stream.LogPrefix = "DB.Backup"
	stream.SinceTs = m.since
	m.since, err = stream.Backup(m.w, m.since)
	return err
}

func (m *restoreOperation) Handle(db *badger.DB) error {
	return db.Load(m.r, 1)
}

func (m *mergeOperation) Handle(db *badger.DB) error {
	const maxBatchCount = 1000
	const maxFrameSize = 256 << 20 // 256 MiB

	br := bufio.NewReaderSize(m.r, 16<<10)
	unmarshalBuf := make([]byte, 1<<10)
	batch := make([]*pb.KV, 0, maxBatchCount)

	flush := func() error {
		if len(batch) == 0 {
			return nil
		}

		for {
			err := db.Update(func(txn *badger.Txn) error {
				for _, kv := range batch {
					if err := m.mergeFn(txn, kv); err != nil {
						return err
					}
				}
				return nil
			})

			switch err {
			case nil:
				batch = batch[:0]
				return nil

			case badger.ErrConflict:
				db.Opts().Logger.Debugf("Conflict detected, retrying merge batch")
				continue

			default:
				return err
			}
		}
	}

	for {
		var sz uint64
		err := binary.Read(br, binary.LittleEndian, &sz)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if sz > maxFrameSize {
			return fmt.Errorf("merge frame size %d exceeds limit %d", sz, maxFrameSize)
		}

		if cap(unmarshalBuf) < int(sz) {
			unmarshalBuf = make([]byte, sz)
		}
		if _, err := io.ReadFull(br, unmarshalBuf[:sz]); err != nil {
			return err
		}

		list := &pb.KVList{}
		if err := proto.Unmarshal(unmarshalBuf[:sz], list); err != nil {
			return err
		}

		for _, kv := range list.Kv {
			batch = append(batch, kv)
			if len(batch) >= maxBatchCount {
				if err := flush(); err != nil {
					return err
				}
			}
		}
	}

	return flush()
}

func (m *viewOperation) Handle(db *badger.DB) error {
	return db.View(m.txnFn)
}

func (m *updateOperation) Handle(db *badger.DB) error {
	for {
		err := db.Update(m.txnFn)

		switch err {
		case nil:
			return nil

		case badger.ErrConflict:
			db.Opts().Logger.Debugf("Conflict detected, retrying transaction")
			continue

		default:
			return err
		}
	}
}

func (m *latestVersionOperation) Handle(db *badger.DB) error {
	m.version = db.MaxVersion()
	return nil
}

func (m *dropAllOperation) Handle(db *badger.DB) error {
	return db.DropAll()
}
