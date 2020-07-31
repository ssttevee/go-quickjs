package internal

// #include "quickjs/quickjs.h"
import "C"

var (
	Null      = Value{tag: C.JS_TAG_NULL}
	Undefined = Value{tag: C.JS_TAG_UNDEFINED}
	False     = Value{tag: C.JS_TAG_BOOL}
	True      = Value{tag: C.JS_TAG_BOOL, u: [8]byte{1, 0, 0, 0, 0, 0, 0, 0}}
)
