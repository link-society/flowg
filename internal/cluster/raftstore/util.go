package raftstore

import "encoding/binary"

func uint64ToBytes(value uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, value)
	return buf
}

func bytesToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}
