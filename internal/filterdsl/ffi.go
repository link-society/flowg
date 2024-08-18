package filterdsl

/*
#cgo LDFLAGS: -L./rust-crate/target/release -lflowg_filterdsl
#include <stdbool.h>
#include <stdlib.h>

typedef struct {
	bool success;
	char *data;
} filterdsl_result;

filterdsl_result filterdsl_compile(const char *input);
void filterdsl_result_free(filterdsl_result result);
*/
import "C"
import "unsafe"

func compile(input string) (string, error) {
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))

	result := C.filterdsl_compile(cInput)
	defer C.filterdsl_result_free(result)

	data := C.GoString(result.data)

	if result.success {
		return data, nil
	} else {
		return "", &CompilationError{Message: data}
	}
}
