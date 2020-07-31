package internal

// #include <stdlib.h>
// #include "quickjs/quickjs.h"
//
// JSValue go_helper_js_throw_syntax_error(JSContext *ctx, const char *str)
// {
//     return JS_ThrowSyntaxError(ctx, str);
// }
//
// JSValue go_helper_js_throw_type_error(JSContext *ctx, const char *str)
// {
//     return JS_ThrowTypeError(ctx, str);
// }
//
// JSValue go_helper_js_throw_reference_error(JSContext *ctx, const char *str)
// {
//     return JS_ThrowReferenceError(ctx, str);
// }
//
// JSValue go_helper_js_throw_range_error(JSContext *ctx, const char *str)
// {
//     return JS_ThrowRangeError(ctx, str);
// }
//
// JSValue go_helper_js_throw_internal_error(JSContext *ctx, const char *str)
// {
//     return JS_ThrowInternalError(ctx, str);
// }
import "C"
import (
	"fmt"
	"unsafe"
)

func ThrowSyntaxError(ctx *Context, format string, v ...interface{}) Value {
	str := C.CString(fmt.Sprintf(format, v...))
	defer C.free(unsafe.Pointer(str))

	return Value(C.go_helper_js_throw_syntax_error((*C.JSContext)(ctx), str))
}

func ThrowTypeError(ctx *Context, format string, v ...interface{}) Value {
	str := C.CString(fmt.Sprintf(format, v...))
	defer C.free(unsafe.Pointer(str))

	return Value(C.go_helper_js_throw_type_error((*C.JSContext)(ctx), str))
}

func ThrowReferenceError(ctx *Context, format string, v ...interface{}) Value {
	str := C.CString(fmt.Sprintf(format, v...))
	defer C.free(unsafe.Pointer(str))

	return Value(C.go_helper_js_throw_reference_error((*C.JSContext)(ctx), str))
}

func ThrowRangeError(ctx *Context, format string, v ...interface{}) Value {
	str := C.CString(fmt.Sprintf(format, v...))
	defer C.free(unsafe.Pointer(str))

	return Value(C.go_helper_js_throw_range_error((*C.JSContext)(ctx), str))
}

func ThrowInternalError(ctx *Context, format string, v ...interface{}) Value {
	str := C.CString(fmt.Sprintf(format, v...))
	defer C.free(unsafe.Pointer(str))

	return Value(C.go_helper_js_throw_internal_error((*C.JSContext)(ctx), str))
}
