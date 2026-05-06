package vrl

/*
#cgo LDFLAGS: -L./rust-crate/target/release -lflowg_vrl -lm
#cgo darwin LDFLAGS: -framework ApplicationServices
#include "ffi.h"
*/
import "C"

import (
	"unsafe"

	"bytes"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

type ScriptRunner struct {
	ffiObject C.script_runner
	buffer    bytes.Buffer
	encoder   *msgpack.Encoder
}

func NewScriptRunner(source string) (*ScriptRunner, error) {
	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))

	res := C.compile_script(cSource)
	if !bool(res.is_ok) {
		defer C.free(unsafe.Pointer(res.err.reason))
		return nil, &CompileError{Message: C.GoString(res.err.reason)}
	}
	obj := res.ok.runner

	self := &ScriptRunner{ffiObject: obj}
	self.buffer.Grow(1024)
	self.encoder = msgpack.NewEncoder(&self.buffer)
	return self, nil
}

func (sr *ScriptRunner) Close() {
	if sr.ffiObject != nil {
		C.drop_script_runner(sr.ffiObject)
		sr.ffiObject = nil
	}
}

func (sr *ScriptRunner) TransformLog(logEvent map[string]string) ([]map[string]string, error) {
	sr.buffer.Reset()

	if err := sr.encoder.Encode(logEvent); err != nil {
		return nil, err
	}
	data := sr.buffer.Bytes()

	buf := C.msgpack_buffer{
		data: (*C.uint8_t)(unsafe.Pointer(&data[0])),
		len:  C.size_t(len(data)),
	}
	res := C.transform_log(sr.ffiObject, buf)
	if !bool(res.is_ok) {
		defer C.free(unsafe.Pointer(res.err.reason))
		return nil, &EvalError{Message: C.GoString(res.err.reason)}
	}

	data = unsafe.Slice((*byte)(unsafe.Pointer(res.ok.data.data)), int(res.ok.data.len))
	var result any
	if err := msgpack.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return normalizeTranformedLog(result)
}

func normalizeTranformedLog(result any) ([]map[string]string, error) {
	switch v := result.(type) {
	case []any:
		logs := make([]map[string]string, len(v))
		for i, item := range v {
			logs[i] = normalizeLogEvent(item)
		}
		return logs, nil

	default:
		return []map[string]string{normalizeLogEvent(result)}, nil
	}
}

func normalizeLogEvent(event any) map[string]string {
	switch v := event.(type) {
	case map[string]any:
		return flattenObjectMap(v)

	case []any:
		return flattenArray(v)

	default:
		return map[string]string{"value": toString(v)}
	}
}

func flattenObjectMap(obj map[string]any) map[string]string {
	flat := make(map[string]string)

	for key, value := range obj {
		switch v := value.(type) {
		case map[string]any:
			nested := flattenObjectMap(v)
			for nestedKey, nestedValue := range nested {
				flatKey := fmt.Sprintf("%s.%s", key, nestedKey)
				flat[flatKey] = nestedValue
			}

		case []any:
			nested := flattenArray(v)
			for nestedKey, nestedValue := range nested {
				flatKey := fmt.Sprintf("%s.%s", key, nestedKey)
				flat[flatKey] = nestedValue
			}

		default:
			flat[key] = toString(v)
		}
	}

	return flat
}

func flattenArray(arr []any) map[string]string {
	flat := make(map[string]string)

	for i, value := range arr {
		key := fmt.Sprintf("%d", i)

		switch v := value.(type) {
		case map[string]any:
			nested := flattenObjectMap(v)
			for nestedKey, nestedValue := range nested {
				flatKey := fmt.Sprintf("%s.%s", key, nestedKey)
				flat[flatKey] = nestedValue
			}

		case []any:
			nested := flattenArray(v)
			for nestedKey, nestedValue := range nested {
				flatKey := fmt.Sprintf("%s.%s", key, nestedKey)
				flat[flatKey] = nestedValue
			}

		default:
			flat[key] = toString(v)
		}
	}

	return flat
}

func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}
