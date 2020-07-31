package internal

// #include "quickjs/quickjs.h"
import "C"

type PropertyFlag int

const (
	PropertyFlagConfigurable    PropertyFlag = C.JS_PROP_CONFIGURABLE
	PropertyFlagWritable        PropertyFlag = C.JS_PROP_WRITABLE
	PropertyFlagEnumerable      PropertyFlag = C.JS_PROP_ENUMERABLE
	PropertyFlagHasGet          PropertyFlag = C.JS_PROP_HAS_GET
	PropertyFlagHasSet          PropertyFlag = C.JS_PROP_HAS_SET
	PropertyFlagHasValue        PropertyFlag = C.JS_PROP_HAS_VALUE
	PropertyFlagHasConfigurable PropertyFlag = C.JS_PROP_HAS_CONFIGURABLE
	PropertyFlagHasWritable     PropertyFlag = C.JS_PROP_HAS_WRITABLE
	PropertyFlagHasEnumerable   PropertyFlag = C.JS_PROP_HAS_ENUMERABLE
	PropertyFlagThrow           PropertyFlag = C.JS_PROP_THROW
	PropertyFlagNoExotic        PropertyFlag = C.JS_PROP_NO_EXOTIC
)

type EvalFlag int

const (
	EvalTypeGlobal           EvalFlag = C.JS_EVAL_TYPE_GLOBAL
	EvalTypeModule           EvalFlag = C.JS_EVAL_TYPE_MODULE
	EvalTypeDirect           EvalFlag = C.JS_EVAL_TYPE_DIRECT
	EvalTypeIndirect         EvalFlag = C.JS_EVAL_TYPE_INDIRECT
	EvalTypeMask             EvalFlag = C.JS_EVAL_TYPE_MASK
	EvalFlagStrict           EvalFlag = C.JS_EVAL_FLAG_STRICT
	EvalFlagStrip            EvalFlag = C.JS_EVAL_FLAG_STRIP
	EvalFlagCompileOnly      EvalFlag = C.JS_EVAL_FLAG_COMPILE_ONLY
	EvalFlagBacktraceBarrier EvalFlag = C.JS_EVAL_FLAG_BACKTRACE_BARRIER
)

type WriteObjectFlag int

const (
	WriteObjectBytecode          WriteObjectFlag = C.JS_WRITE_OBJ_BYTECODE
	WriteObjectByteSwapped       WriteObjectFlag = C.JS_WRITE_OBJ_BSWAP
	WriteObjectSharedArrayBuffer WriteObjectFlag = C.JS_WRITE_OBJ_SAB
	WriteObjectReference         WriteObjectFlag = C.JS_WRITE_OBJ_REFERENCE
)

type ReadObjectFlag int

const (
	ReadObjectBytecode          ReadObjectFlag = C.JS_READ_OBJ_BYTECODE
	ReadObjectNoCopy            ReadObjectFlag = C.JS_READ_OBJ_ROM_DATA
	ReadObjectSharedArrayBuffer ReadObjectFlag = C.JS_READ_OBJ_SAB
	ReadObjectReference         ReadObjectFlag = C.JS_READ_OBJ_REFERENCE
)
