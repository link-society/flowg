package sse

import "bytes"

func findDoubleNewLine(data []byte) (int, int) {
	crcr := bytes.Index(data, []byte("\r\r"))
	lflf := bytes.Index(data, []byte("\n\n"))
	crlflf := bytes.Index(data, []byte("\r\n\n"))
	lfcrlf := bytes.Index(data, []byte("\n\r\n"))
	crlfcrlf := bytes.Index(data, []byte("\r\n\r\n"))

	minPos := minPosInts(crcr, lflf, crlflf, lfcrlf, crlfcrlf)

	nlen := 2
	if minPos == crlfcrlf {
		nlen = 4
	} else if minPos == crlflf || minPos == lfcrlf {
		nlen = 3
	}

	return minPos, nlen
}

func minPosInts(value int, values ...int) int {
	var min int = value
	for _, v := range values {
		min = minPosInt(min, v)
	}
	return min
}

func minPosInt(a, b int) int {
	if a < 0 {
		return b
	}
	if b < 0 {
		return a
	}
	if a > b {
		return b
	}
	return a
}
