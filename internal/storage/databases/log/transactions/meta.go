package transactions

import (
	"fmt"

	"encoding/json"

	"link-society.com/flowg/internal/app/featureflags"

	"link-society.com/flowg/internal/models"
	"link-society.com/flowg/internal/storage/generic/kv"
)

var demoStreamConfig = models.StreamConfig{
	RetentionTime: 15 * 60, // 15 minutes
	RetentionSize: 10,
	IndexedFields: []string{},
}

// FetchStreamConfigs returns the configuration of every stream, keyed by stream
// name. In demo mode the fixed demo configuration is substituted for all of
// them.
func FetchStreamConfigs(txn kv.QueryTx) (map[string]models.StreamConfig, error) {
	streams := map[string]models.StreamConfig{}

	for pair := range txn.IterPairs(kv.Key{"stream", "config"}, kv.KeyRange{}) {
		key := pair.Key()

		stream := key[len(key)-1]

		var streamConfig models.StreamConfig
		if featureflags.GetDemoMode() {
			streamConfig = demoStreamConfig
		} else {
			val := pair.Value()
			// If the value is empty, we treat it as a default config (the zero-value
			// of StreamConfig).
			if len(val) > 0 {
				if err := json.Unmarshal(val, &streamConfig); err != nil {
					return nil, fmt.Errorf("could not unmarshal stream config '%s': %w", stream, err)
				}
			}
		}

		if streamConfig.IndexedFields == nil {
			streamConfig.IndexedFields = []string{}
		}

		streams[stream] = streamConfig
	}

	return streams, nil
}

// FetchStreamFields lists the field names recorded for a stream from its
// "stream:field:<stream>:" existence markers.
func FetchStreamFields(txn kv.QueryTx, stream string) []string {
	fields := []string{}

	for key := range txn.IterKeys(kv.Key{"stream", "field", stream}, kv.KeyRange{}) {
		fieldName := key[len(key)-1]
		fields = append(fields, fieldName)
	}

	return fields
}

// GetOrCreateStreamConfig returns a stream's configuration, lazily creating an
// empty (default) config the first time a stream is referenced so later writes
// have something to update. Demo mode short-circuits to the fixed config.
func GetOrCreateStreamConfig(txn kv.MutationTx, stream string) (models.StreamConfig, error) {
	if featureflags.GetDemoMode() {
		return demoStreamConfig, nil
	}

	var streamConfig models.StreamConfig

	streamKey := kv.Key{"stream", "config", stream}
	val, err := txn.Get(streamKey)
	if err != nil {
		return models.StreamConfig{}, fmt.Errorf(
			"could not fetch stream config '%s': %w",
			stream, err,
		)
	}

	if val == nil {
		err := txn.Set(streamKey, []byte{})
		if err != nil {
			return models.StreamConfig{}, fmt.Errorf(
				"could not create default stream config '%s': %w",
				stream, err,
			)
		}
	}

	// If the value is empty, we treat it as a default config (the zero-value
	// of StreamConfig).
	if len(val) > 0 {
		if err := json.Unmarshal(val, &streamConfig); err != nil {
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

// ConfigureStream stores a new stream configuration and reconciles its indexed
// fields against the previous one: fields newly added to the index are back-
// filled over existing records, and fields removed from it have their index
// dropped.
func ConfigureStream(txn kv.MutationTx, stream string, config models.StreamConfig) error {
	if config.IndexedFields == nil {
		config.IndexedFields = []string{}
	}

	oldConfig, err := GetOrCreateStreamConfig(txn, stream)
	if err != nil {
		return fmt.Errorf("could not fetch old stream config '%s': %w", stream, err)
	}

	streamKey := kv.Key{"stream", "config", stream}
	configVal, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("could not marshal stream config '%s': %w", stream, err)
	}

	if err := txn.Set(streamKey, configVal); err != nil {
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

// DeleteStream removes everything belonging to a stream: its log entries
// ("entry:<stream>:*"), inverted indexes ("index:<stream>:*"), field markers
// ("stream:field:<stream>:*") and its configuration key.
func DeleteStream(txn kv.MutationTx, stream string) error {
	prefixes := []kv.Key{
		{"entry", stream},
		{"index", stream},
		{"stream", "field", stream},
	}

	for _, prefix := range prefixes {
		for key := range txn.IterKeys(prefix, kv.KeyRange{}) {
			if err := txn.Clear(key); err != nil {
				return fmt.Errorf(
					"could not delete key '%s' from stream '%s': %w",
					key, stream, err,
				)
			}
		}
	}

	streamKey := kv.Key{"stream", "config", stream}
	if err := txn.Clear(streamKey); err != nil {
		return fmt.Errorf("could not delete stream config '%s': %w", stream, err)
	}

	return nil
}
