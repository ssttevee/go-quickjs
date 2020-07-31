package internal

// #include "quickjs/quickjs.h"
import "C"

type Tag int64

const (
	TagBigDecimal       Tag = C.JS_TAG_BIG_DECIMAL
	TagBigInt           Tag = C.JS_TAG_BIG_INT
	TagBigFloat         Tag = C.JS_TAG_BIG_FLOAT
	TagSymbol           Tag = C.JS_TAG_SYMBOL
	TagString           Tag = C.JS_TAG_STRING
	TagModule           Tag = C.JS_TAG_MODULE
	TagFunctionBytecode Tag = C.JS_TAG_FUNCTION_BYTECODE
	TagObject           Tag = C.JS_TAG_OBJECT

	TagInt           Tag = C.JS_TAG_INT
	TagBool          Tag = C.JS_TAG_BOOL
	TagNull          Tag = C.JS_TAG_NULL
	TagUndefined     Tag = C.JS_TAG_UNDEFINED
	TagUninitialized Tag = C.JS_TAG_UNINITIALIZED
	TagCatchOffset   Tag = C.JS_TAG_CATCH_OFFSET
	TagException     Tag = C.JS_TAG_EXCEPTION
	TagFloat64       Tag = C.JS_TAG_FLOAT64
)
