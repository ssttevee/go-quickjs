package internal

// #include <stdlib.h>
// #include "quickjs/quickjs.h"
import "C"
import (
	"reflect"
	"unsafe"
)

func makeSliceHeader(data unsafe.Pointer, n int) unsafe.Pointer {
	return unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(data)), Len: n, Cap: n})
}

func jsvalues(args []Value) []C.JSValue {
	converted := make([]C.JSValue, len(args))
	for i, arg := range args {
		converted[i] = C.JSValue(arg)
	}

	return converted
}

func argv(args []C.JSValue) *C.JSValue {
	return (*C.JSValue)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&args)).Data))
}

func booltoint(b bool) int {
	if b {
		return 1
	}

	return 0
}
