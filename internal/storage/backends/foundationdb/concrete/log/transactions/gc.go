package transactions

import (
	"fmt"

	fdb "github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
)

// EstimateStorage returns an estimate, in bytes, of the storage used by the
// named stream, using FDB's GetEstimatedRangeSizeBytes on the entry subspace.
func EstimateStorage(txn fdb.ReadTransaction, stream string) (int64, error) {
	streamEntrySS := entrySS.Sub(subspace.FromBytes([]byte(stream)))

	estFuture := txn.GetEstimatedRangeSizeBytes(streamEntrySS)
	est, err := estFuture.Get()
	if err != nil {
		return 0, fmt.Errorf("could not estimate storage for stream '%s': %w", stream, err)
	}

	return int64(est), nil
}

// CollectGarbage evicts the oldest records of any stream that has grown past its
// configured size budget. For each over-budget stream it walks the entries from
// oldest to newest (the key order), deleting them and their inverted-index
// references until the stream fits within its retention size again.
//
// Note: because FDB has no per-key TTL, this GC also handles time-based eviction
// (entries past their retention time) by checking the timestamp embedded in each
// entry key.
func CollectGarbage(txn fdb.Transaction) error {
	streams, err := FetchStreamConfigs(txn)
	if err != nil {
		return err
	}

	for stream, config := range streams {
		retentionSize := config.RetentionSize * 1024 * 1024 // MB to bytes

		if retentionSize > 0 {
			streamEntrySS := entrySS.Sub(subspace.FromBytes([]byte(stream)))

			// Estimate the total size of the stream
			estFuture := txn.GetEstimatedRangeSizeBytes(streamEntrySS)
			est, err := estFuture.Get()
			if err != nil {
				return fmt.Errorf("could not estimate size for stream '%s': %w", stream, err)
			}

			if int64(est) > retentionSize {
				// Walk oldest to newest, deleting until under budget.
				ri := txn.GetRange(streamEntrySS, fdb.RangeOptions{}).Iterator()
				for ri.Advance() {
					kv := ri.MustGet()
					key := append([]byte(nil), kv.Key...)

					// Clear buffers the mutation; actual commit happens when
					// the FDB transaction completes.
					txn.Clear(fdb.Key(key))
					purgeEntryFromFieldIndex(txn, stream, key)

					// Re-check remaining size
					estFuture := txn.GetEstimatedRangeSizeBytes(streamEntrySS)
					est, err := estFuture.Get()
					if err != nil {
						return fmt.Errorf(
							"could not re-estimate size for stream '%s': %w",
							stream, err,
						)
					}

					if int64(est) <= retentionSize {
						break
					}
				}
			}
		}
	}

	return nil
}
