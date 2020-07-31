package internal

// #include "quickjs/quickjs.h"
import "C"
import (
	"reflect"
	"unsafe"
)

type PropertyEnum C.JSPropertyEnum

func (p PropertyEnum) Atom() Atom {
	return Atom(p.atom)
}

func (p PropertyEnum) IsEnumerable() bool {
	return p.is_enumerable != 0
}

func GetOwnPropertyNames(ctx *Context, obj Value, flags int) []PropertyEnum {
	var tab *C.JSPropertyEnum
	var size uint32
	if C.JS_GetOwnPropertyNames((*C.JSContext)(ctx), (**C.JSPropertyEnum)(unsafe.Pointer(&tab)), (*C.uint32_t)(&size), C.JSValue(obj), C.int(flags)) < 0 {
		return nil
	}

	return *(*[]PropertyEnum)(makeSliceHeader(unsafe.Pointer(tab), int(size)))
}

func FreePropertyEnum(ctx *Context, tab []PropertyEnum) {
	if tab == nil {
		return
	}

	for _, p := range tab {
		C.JS_FreeAtom((*C.JSContext)(ctx), p.atom)
	}

	C.js_free((*C.JSContext)(ctx), unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&tab)).Data))
}
