package transactions

import (
	"github.com/apple/foundationdb/bindings/go/src/fdb/subspace"
)

// Top-level subspace prefixes for the FDB log storage key space.
//
// The hierarchy is:
//
//	entry        -> entrySS
//	  <stream>   -> entrySS.Sub(stream)       (raw prefix)
//	    <ts>:<uuid>  -> streamEntrySS.Pack({ts, uuid})
//
//	config       -> cfgSS
//	  <stream>   -> cfgSS.Pack({stream})       (tuple element)
//
//	field        -> fieldSS
//	  <stream>   -> fieldSS.Sub(stream)        (raw prefix)
//	    <name>   -> streamFieldSS.Pack({name}) (tuple element)
//
//	index        -> indexSS
//	  <stream>   -> indexSS.Sub(stream)        (raw prefix)
//	    <field>  -> streamIndexSS.Sub(field)   (raw prefix)
//	      <base64_value>
//	               -> fieldIndexSS.Sub(b64)    (raw prefix)
//	        <entry_key>
//	               -> valueIndexSS.Pack({entryKey})
var (
	entrySS = subspace.FromBytes([]byte("entry"))
	cfgSS   = subspace.FromBytes([]byte("config"))
	fieldSS = subspace.FromBytes([]byte("field"))
	indexSS = subspace.FromBytes([]byte("index"))
)
