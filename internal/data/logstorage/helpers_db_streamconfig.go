package logstorage

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v3"
)

func getOrCreateStreamConfig(txn *badger.Txn, stream string) (StreamConfig, error) {
	var streamConfig StreamConfig

	streamKey := []byte(fmt.Sprintf("stream:%s", stream))
	switch streamConfigItem, err := txn.Get(streamKey); {
	case err != nil && err != badger.ErrKeyNotFound:
		return StreamConfig{}, fmt.Errorf(
			"could not fetch stream config '%s': %w",
			stream, err,
		)

	case err == badger.ErrKeyNotFound:
		err := txn.Set(streamKey, []byte(""))
		if err != nil {
			return StreamConfig{}, fmt.Errorf(
				"could not create default stream config '%s': %w",
				stream, err,
			)
		}

	case err == nil:
		err := streamConfigItem.Value(func(val []byte) error {
			if len(val) > 0 {
				if err := json.Unmarshal(val, &streamConfig); err != nil {
					return fmt.Errorf(
						"could not unmarshal stream config '%s': %w",
						stream, err,
					)
				}
			}
			return nil
		})
		if err != nil {
			return StreamConfig{}, err
		}
	}

	return streamConfig, nil
}

func fetchStreamConfigs(txn *badger.Txn) (map[string]StreamConfig, error) {
	streams := map[string]StreamConfig{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = []byte("stream:")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		stream := string(it.Item().Key()[7:])

		var streamConfig StreamConfig
		err := it.Item().Value(func(val []byte) error {
			if len(val) > 0 {
				if err := json.Unmarshal(val, &streamConfig); err != nil {
					return fmt.Errorf("could not unmarshal stream config '%s': %w", stream, err)
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}

		streams[stream] = streamConfig
	}

	return streams, nil
}
