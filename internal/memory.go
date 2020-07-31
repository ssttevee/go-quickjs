package internal

// #include <stdlib.h>
// #include "quickjs/quickjs.h"
//
// extern void *go_custom_js_malloc(JSMallocState *s, size_t size);
// extern void go_custom_js_free(JSMallocState *s, void *ptr);
// extern void *go_custom_js_realloc(JSMallocState *s, void *ptr, size_t size);
import "C"
import (
	"unsafe"
)

var (
	// use cgo memory functions to guarantee allocations never returns nil
	gomallocfuncs = C.JSMallocFunctions{
		js_malloc:  (*[0]byte)(C.go_custom_js_malloc),
		js_free:    (*[0]byte)(C.go_custom_js_free),
		js_realloc: (*[0]byte)(C.go_custom_js_realloc),
	}
)

//export go_custom_js_malloc
func go_custom_js_malloc(s *C.JSMallocState, size csize) unsafe.Pointer {
	return C.malloc(size)
}

//export go_custom_js_free
func go_custom_js_free(s *C.JSMallocState, ptr unsafe.Pointer) {
	C.free(ptr)
}

//export go_custom_js_realloc
func go_custom_js_realloc(s *C.JSMallocState, ptr unsafe.Pointer, size csize) unsafe.Pointer {
	return C.realloc(ptr, size)
}
