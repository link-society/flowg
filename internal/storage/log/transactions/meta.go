package transactions

import (
	"fmt"

	"encoding/json"

	"github.com/dgraph-io/badger/v4"

	"link-society.com/flowg/internal/app/featureflags"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/utils/hlc"
	"link-society.com/flowg/internal/utils/lww"
)

var demoStreamConfig = models.StreamConfig{
	RetentionTime: 15 * 60, // 15 minutes
	RetentionSize: 10,
	IndexedFields: []string{},
}

func FetchStreamConfigs(txn *badger.Txn) (map[string]models.StreamConfig, error) {
	streams := map[string]models.StreamConfig{}

	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = true
	opts.Prefix = []byte("stream:config:")
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		stream := string(it.Item().Key()[len(opts.Prefix):])

		var streamConfig models.StreamConfig
		if featureflags.GetDemoMode() {
			streamConfig = demoStreamConfig
		} else {
			live := false
			err := it.Item().Value(func(val []byte) error {
				env, err := lww.Unmarshal(val)
				if err != nil {
					return fmt.Errorf("could not unmarshal stream config '%s': %w", stream, err)
				}
				if env.Deleted {
					return nil
				}
				live = true
				if len(env.Payload) > 0 {
					if err := json.Unmarshal(env.Payload, &streamConfig); err != nil {
						return fmt.Errorf("could not unmarshal stream config '%s': %w", stream, err)
					}
				}
				return nil
			})
			if err != nil {
				return nil, err
			}
			if !live {
				continue
			}
		}

		if streamConfig.IndexedFields == nil {
			streamConfig.IndexedFields = []string{}
		}

		streams[stream] = streamConfig
	}

	return streams, nil
}

func FetchStreamFields(txn *badger.Txn, stream string) []string {
	fields := []string{}

	prefix := []byte(fmt.Sprintf("stream:field:%s:", stream))
	opts := badger.DefaultIteratorOptions
	opts.PrefetchValues = false
	opts.Prefix = prefix
	it := txn.NewIterator(opts)
	defer it.Close()

	for it.Rewind(); it.Valid(); it.Next() {
		item := it.Item()
		key := item.Key()

		fieldName := string(key[len(prefix):])
		fields = append(fields, fieldName)
	}

	return fields
}

func GetOrCreateStreamConfig(txn *badger.Txn, stream string, ts hlc.Timestamp) (models.StreamConfig, error) {
	if featureflags.GetDemoMode() {
		return demoStreamConfig, nil
	}

	var streamConfig models.StreamConfig

	streamKey := []byte(fmt.Sprintf("stream:config:%s", stream))
	env, found, err := lww.Read(txn, streamKey)
	if err != nil {
		return models.StreamConfig{}, fmt.Errorf(
			"could not fetch stream config '%s': %w",
			stream, err,
		)
	}

	if !found {
		if _, err := lww.Apply(txn, streamKey, lww.Envelope{Timestamp: ts}); err != nil {
			return models.StreamConfig{}, fmt.Errorf(
				"could not create default stream config '%s': %w",
				stream, err,
			)
		}
	} else if len(env.Payload) > 0 {
		if err := json.Unmarshal(env.Payload, &streamConfig); err != nil {
			return models.StreamConfig{}, fmt.Errorf(
				"could not unmarshal stream config '%s': %w",
				stream, err,
			)
		}
	}

	if streamConfig.IndexedFields == nil {
		streamConfig.IndexedFields = []string{}
	}

	return streamConfig, nil
}

func ConfigureStream(txn *badger.Txn, stream string, config models.StreamConfig, ts hlc.Timestamp) error {
	if config.IndexedFields == nil {
		config.IndexedFields = []string{}
	}

	oldConfig, err := GetOrCreateStreamConfig(txn, stream, ts)
	if err != nil {
		return fmt.Errorf("could not fetch old stream config '%s': %w", stream, err)
	}

	streamKey := []byte(fmt.Sprintf("stream:config:%s", stream))
	configVal, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal stream config '%s': %w", stream, err)
	}

	if _, err := lww.Apply(txn, streamKey, lww.Envelope{Timestamp: ts, Payload: configVal}); err != nil {
		return fmt.Errorf("could not save stream config '%s': %w", stream, err)
	}

	for _, field := range config.IndexedFields {
		if !oldConfig.IsFieldIndexed(field) {
			IndexField(txn, stream, field)
		}
	}

	for _, field := range oldConfig.IndexedFields {
		if !config.IsFieldIndexed(field) {
			UnindexField(txn, stream, field)
		}
	}

	return nil
}

func DeleteStream(txn *badger.Txn, stream string, ts hlc.Timestamp) error {
	prefixes := []string{
		fmt.Sprintf("entry:%s:", stream),
		fmt.Sprintf("index:%s:", stream),
		fmt.Sprintf("stream:field:%s:", stream),
	}

	for _, prefix := range prefixes {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().KeyCopy(nil)

			if err := txn.Delete(key); err != nil {
				return fmt.Errorf(
					"could not delete key '%s' from stream '%s': %w",
					key, stream, err,
				)
			}
		}
	}

	streamKey := []byte(fmt.Sprintf("stream:config:%s", stream))
	if _, err := lww.Apply(txn, streamKey, lww.Envelope{Timestamp: ts, Deleted: true}); err != nil {
		return fmt.Errorf("could not delete stream config '%s': %w", stream, err)
	}

	return nil
}
