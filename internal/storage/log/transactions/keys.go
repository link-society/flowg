package transactions

import (
	"bytes"
	"strings"
)

var (
	entryKeyPrefix  = []byte("entry:")
	configKeyPrefix = []byte("stream:config:")
)

func StreamFromEntryKey(key []byte) (string, bool) {
	if !bytes.HasPrefix(key, entryKeyPrefix) {
		return "", false
	}

	rest := string(key[len(entryKeyPrefix):])

	uuidSep := strings.LastIndexByte(rest, ':')
	if uuidSep < 0 {
		return "", false
	}

	tsSep := strings.LastIndexByte(rest[:uuidSep], ':')
	if tsSep < 0 {
		return "", false
	}

	return rest[:tsSep], true
}

func StreamFromConfigKey(key []byte) (string, bool) {
	if !bytes.HasPrefix(key, configKeyPrefix) {
		return "", false
	}
	return string(key[len(configKeyPrefix):]), true
}
