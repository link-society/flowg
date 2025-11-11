package vrl

/*
#cgo LDFLAGS: -L./rust-crate/target/release -lflowg_vrl -lm
#include "ffi.h"
*/
import "C"

import (
	"unsafe"
)

type ScriptRunner struct {
	ffiObject C.script_runner
}

func NewScriptRunner(source string) (*ScriptRunner, error) {
	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))

	var err *C.char
	obj := C.script_runner_new(cSource, &err)
	if obj == nil {
		defer C.free(unsafe.Pointer(err))
		return nil, &CompileError{Message: C.GoString(err)}
	}

	return &ScriptRunner{ffiObject: obj}, nil
}

func (sr *ScriptRunner) Close() {
	if sr.ffiObject != nil {
		C.script_runner_free(sr.ffiObject)
		sr.ffiObject = nil
	}
}

func (sr *ScriptRunner) Eval(input map[string]string) (map[string]string, error) {
	cInput := mapToHmap(input)
	defer freeHmap(cInput)

	var err *C.char
	cOutput := C.script_runner_eval(sr.ffiObject, cInput, &err)
	defer C.hmap_free(cOutput)

	if cOutput == nil {
		defer C.free(unsafe.Pointer(err))
		return nil, &EvalError{Message: C.GoString(err)}
	}

	return hmapToMap(cOutput), nil
}

func hmapToMap(cHmap *C.hmap) map[string]string {
	result := make(map[string]string)

	if cHmap != nil && cHmap.count > 0 {
		count := int(cHmap.count)
		cEntries := unsafe.Pointer(cHmap.entries)

		if cEntries != nil {
			cSlice := (*[1 << 30]C.hmap_entry)(cEntries)[:count:count]
			for _, entry := range cSlice {
				key := C.GoString(entry.key)
				value := C.GoString(entry.value)
				result[key] = value
			}
		}
	}

	return result
}

func mapToHmap(input map[string]string) *C.hmap {
	var (
		cEntries    unsafe.Pointer
		cEntrySlice []C.hmap_entry
	)

	if len(input) > 0 {
		cEntries = C.malloc(C.size_t(len(input)) * C.size_t(unsafe.Sizeof(C.hmap_entry{})))

		cEntrySlice = (*[1 << 30]C.hmap_entry)(cEntries)[:len(input):len(input)]
		i := 0
		for key, value := range input {
			cKey := C.CString(key)
			cValue := C.CString(value)
			cEntrySlice[i] = C.hmap_entry{key: cKey, value: cValue}
			i++
		}
	}

	cHmap := (*C.hmap)(C.malloc(C.size_t(unsafe.Sizeof(C.hmap{}))))
	cHmap.count = C.size_t(len(input))
	cHmap.entries = (*C.hmap_entry)(cEntries)

	return cHmap
}

func freeHmap(cHmap *C.hmap) {
	if cHmap != nil {
		count := int(cHmap.count)
		cEntries := unsafe.Pointer(cHmap.entries)

		if cEntries != nil {
			cSlice := (*[1 << 30]C.hmap_entry)(cEntries)[:count:count]
			for _, entry := range cSlice {
				C.free(unsafe.Pointer(entry.key))
				C.free(unsafe.Pointer(entry.value))
			}
			C.free(cEntries)
		}
		C.free(unsafe.Pointer(cHmap))
	}
}
