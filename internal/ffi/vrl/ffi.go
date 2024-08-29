package vrl

/*
#cgo LDFLAGS: -L./rust-crate/target/release -lflowg_vrl
#include <stdlib.h>

typedef struct {
	char* key;
	char* value;
} hmap_entry;

typedef struct {
	size_t count;
	hmap_entry* entries;
} hmap;

typedef enum {
	vrl_result_ok,
	vrl_result_err
} vrl_result_tag;

typedef struct {
	hmap* ok_data;
	char* err_data;
} vrl_result_data;

typedef struct {
	vrl_result_tag tag;
	vrl_result_data data;
} vrl_result;

extern vrl_result* process_record(hmap* input, const char* script);
extern void vrl_result_free(vrl_result* result);
*/
import "C"

import (
	"unsafe"
)

func ProcessRecord(
	record map[string]string,
	script string,
) (map[string]string, error) {
	cEntries := C.malloc(C.size_t(len(record)) * C.size_t(unsafe.Sizeof(C.hmap_entry{})))
	defer C.free(cEntries)

	cEntrySlice := (*[1 << 30]C.hmap_entry)(cEntries)[:len(record):len(record)]
	i := 0
	for key, value := range record {
		cKey := C.CString(key)
		cValue := C.CString(value)
		cEntrySlice[i] = C.hmap_entry{key: cKey, value: cValue}
		i++
	}

	cInput := (*C.hmap)(C.malloc(C.size_t(unsafe.Sizeof(C.hmap{}))))
	defer C.free(unsafe.Pointer(cInput))

	cInput.count = C.size_t(len(record))
	cInput.entries = (*C.hmap_entry)(cEntries)

	cScript := C.CString(script)
	defer C.free(unsafe.Pointer(cScript))

	cResult := C.process_record(cInput, cScript)
	defer C.vrl_result_free(cResult)

	result := make(map[string]string)
	if cResult == nil {
		return nil, &NullPointerError{}
	}

	switch cResult.tag {
	case C.vrl_result_ok:
		cResultSlice := (*[1 << 30]C.hmap_entry)(unsafe.Pointer(cResult.data.ok_data.entries))[:cResult.data.ok_data.count:cResult.data.ok_data.count]
		for _, entry := range cResultSlice {
			key := C.GoString(entry.key)
			value := C.GoString(entry.value)
			result[key] = value
		}
		return result, nil

	case C.vrl_result_err:
		err := C.GoString(cResult.data.err_data)
		return nil, &RuntimeError{Message: err}

	default:
		return nil, &RuntimeError{Message: "unknown error"}
	}
}
