package transactions

import (
	"fmt"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
	"github.com/apple/foundationdb/bindings/go/src/fdb/tuple"
)

// Top-level subspace prefixes for the FDB log storage key space.
//
// The hierarchy is:
//
//	prefix (optional root namespace, e.g. "log/")
//	  entry        -> entrySS
//	    <stream>   -> entrySS.Sub(stream)       (raw prefix)
//	      <ts>:<uuid>  -> streamEntrySS.Pack({ts, uuid})
//
//	  config       -> cfgSS
//	    <stream>   -> cfgSS.Pack({stream})       (tuple element)
//
//	  field        -> fieldSS
//	    <stream>   -> fieldSS.Sub(stream)        (raw prefix)
//	      <name>   -> streamFieldSS.Pack({name}) (tuple element)
//
//	  index        -> indexSS
//	    <stream>   -> indexSS.Sub(stream)        (tuple-encoded byte string)
//	      <field>  -> streamIndexSS.Sub(field)   (tuple-encoded byte string)
//	        <base64_value>
//	                 -> fieldIndexSS.Sub(b64)    (tuple-encoded byte string)
//	          <entry_key>
//	                 -> valueIndexSS.Pack({entryKey})
//
// Each Sub() call packs the element as a tuple-encoded byte string:
// 0x01 <bytes> 0x00 (type code, raw bytes, null terminator).
// Default (no-prefix) subspaces for backwards compatibility when Init is
// called with a nil/empty root.
var (
	entrySS = subspace.FromBytes([]byte("entry"))
	cfgSS   = subspace.FromBytes([]byte("config"))
	fieldSS = subspace.FromBytes([]byte("field"))
	indexSS = subspace.FromBytes([]byte("index"))
)

// Init roots all package-level subspaces under the given prefix so that log
// data is isolated from other storage domains. When root is nil or empty the
// default raw subspaces are kept (backwards compatible).
func Init(root subspace.Subspace) {
	if root == nil || len(root.Bytes()) == 0 {
		return
	}
	entrySS = root.Sub(subspace.FromBytes([]byte("entry")))
	cfgSS = root.Sub(subspace.FromBytes([]byte("config")))
	fieldSS = root.Sub(subspace.FromBytes([]byte("field")))
	indexSS = root.Sub(subspace.FromBytes([]byte("index")))
}

// PackEntryKey builds the FDB-packed entry key for a (stream, timestamp, uuid).
// The key sorts by timestamp within each stream.
func PackEntryKey(stream string, ts time.Time, uuidStr string) []byte {
	streamEntrySS := entrySS.Sub(subspace.FromBytes([]byte(stream)))
	paddedTs := fmt.Sprintf("%020d", ts.UnixMilli())
	return streamEntrySS.Pack(tuple.Tuple{paddedTs, uuidStr})
}
