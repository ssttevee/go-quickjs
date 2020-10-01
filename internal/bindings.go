package internal

// The goal of this package is to expose the C functions without any of the C
// types.
//
// Additionally, all of the exposed functions should be normalized so that
// calling `DupValue` is never required before calling.

// #cgo LDFLAGS: -lquickjs -lm -lpthread
// #include <stdlib.h>
// #include "quickjs/quickjs.h"
//
// extern JSValue js_call_func_job(JSContext *ctx, int argc, JSValueConst *argv);
import "C"
import (
	"reflect"
	"runtime"
	"unsafe"
)

type Runtime C.JSRuntime

type Context C.JSContext

type Atom C.JSAtom

type Value C.JSValue

func (v Value) Tag() Tag {
	return Tag(v.tag)
}

func Eval(ctx *Context, input string, filename string, flags EvalFlag) Value {
	cinput := C.CString(input)
	defer C.free(unsafe.Pointer(cinput))

	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	return Value(C.JS_Eval((*C.JSContext)(ctx), cinput, csize(len(input)), cfilename, C.int(flags)))
}

func ToString(ctx *Context, v Value) string {
	var size csize
	cstr := C.JS_ToCStringLen((*C.JSContext)(ctx), &size, C.JSValue(v))
	defer C.JS_FreeCString((*C.JSContext)(ctx), cstr)

	return C.GoStringN(cstr, C.int(size))
}

func ToInt(v Value) int {
	return int(*(*C.int32_t)(unsafe.Pointer(&v.u)))
}

func ToFloat(v Value) float64 {
	return float64(*(*C.double)(unsafe.Pointer(&v.u)))
}

func ToBool(v Value) bool {
	return ToInt(v) != 0
}

func IsTruthy(ctx *Context, v Value) int {
	return int(C.JS_ToBool((*C.JSContext)(ctx), C.JSValue(v)))
}

func GetProperty(ctx *Context, thisObj Value, prop Atom) Value {
	return Value(C.JS_GetProperty((*C.JSContext)(ctx), C.JSValue(thisObj), C.JSAtom(prop)))
}

func GetPropertyStr(ctx *Context, thisObj Value, prop string) Value {
	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))

	return Value(C.JS_GetPropertyStr((*C.JSContext)(ctx), C.JSValue(thisObj), cprop))
}

func GetPropertyInt(ctx *Context, thisObj Value, index int) Value {
	return Value(C.JS_GetPropertyUint32((*C.JSContext)(ctx), C.JSValue(thisObj), C.uint32_t(index)))
}

func SetProperty(ctx *Context, thisObj Value, prop Atom, val Value) int {
	return int(C.JS_SetProperty((*C.JSContext)(ctx), C.JSValue(thisObj), C.JSAtom(prop), C.JSValue(DupValue(ctx, val))))
}

func SetPropertyStr(ctx *Context, thisObj Value, prop string, val Value) int {
	cprop := C.CString(prop)
	defer C.free(unsafe.Pointer(cprop))

	return int(C.JS_SetPropertyStr((*C.JSContext)(ctx), C.JSValue(thisObj), cprop, C.JSValue(DupValue(ctx, val))))
}

func SetPropertyInt(ctx *Context, thisObj Value, index int, val Value) int {
	return int(C.JS_SetPropertyUint32((*C.JSContext)(ctx), C.JSValue(thisObj), C.uint32_t(index), C.JSValue(DupValue(ctx, val))))
}

func Call(ctx *Context, funcObj, thisObj Value, args []Value) Value {
	jsargs := jsvalues(args)
	defer runtime.KeepAlive(jsargs)

	return Value(C.JS_Call((*C.JSContext)(ctx), C.JSValue(funcObj), C.JSValue(thisObj), C.int(len(args)), argv(jsargs)))
}

func InvokeStr(ctx *Context, thisVal Value, name string, args []Value) Value {
	jsargs := jsvalues(args)
	defer runtime.KeepAlive(jsargs)

	atom := NewAtom(ctx, name)
	defer FreeAtom(ctx, atom)

	return Value(C.JS_Invoke((*C.JSContext)(ctx), C.JSValue(thisVal), C.JSAtom(atom), C.int(len(args)), argv(jsargs)))
}

func CallConstructor(ctx *Context, funcObj Value, args []Value) Value {
	jsargs := jsvalues(args)
	defer runtime.KeepAlive(jsargs)

	return Value(C.JS_CallConstructor((*C.JSContext)(ctx), C.JSValue(funcObj), C.int(len(args)), argv(jsargs)))
}

func DefineProperty(ctx *Context, obj Value, prop Atom, v, getter, setter Value, flags PropertyFlag) int {
	return int(C.JS_DefineProperty((*C.JSContext)(ctx), C.JSValue(obj), C.JSAtom(prop), C.JSValue(v), C.JSValue(getter), C.JSValue(setter), C.int(flags)))
}

func DefinePropertyValueStr(ctx *Context, obj Value, prop string, v Value, flags PropertyFlag) int {
	s := C.CString(prop)
	defer C.free(unsafe.Pointer(s))

	return int(C.JS_DefinePropertyValueStr((*C.JSContext)(ctx), C.JSValue(obj), s, C.JSValue(DupValue(ctx, v)), C.int(flags)))
}

func DefinePropertyValue(ctx *Context, obj Value, atom Atom, v Value, flags PropertyFlag) int {
	return int(C.JS_DefinePropertyValue((*C.JSContext)(ctx), C.JSValue(obj), C.JSAtom(atom), C.JSValue(DupValue(ctx, v)), C.int(flags)))
}

func DefinePropertyGetSet(ctx *Context, obj Value, atom Atom, getter, setter Value, flags PropertyFlag) int {
	return int(C.JS_DefinePropertyGetSet((*C.JSContext)(ctx), C.JSValue(obj), C.JSAtom(atom), C.JSValue(DupValue(ctx, getter)), C.JSValue(DupValue(ctx, setter)), C.int(flags)))
}

func NewRuntime() *Runtime {
	return (*Runtime)(defineCustomClasses(C.JS_NewRuntime2(&gomallocfuncs, nil)))
}

func FreeRuntime(rt *Runtime) {
	C.JS_FreeRuntime((*C.JSRuntime)(rt))
}

func NewContext(rt *Runtime) *Context {
	return (*Context)(C.JS_NewContext((*C.JSRuntime)(rt)))
}

func FreeContext(ctx *Context) {
	C.JS_FreeContext((*C.JSContext)(ctx))
}

func GetGlobalObject(ctx *Context) Value {
	return Value(C.JS_GetGlobalObject((*C.JSContext)(ctx)))
}

func NewError(ctx *Context) Value {
	return Value(C.JS_NewError((*C.JSContext)(ctx)))
}

func DupValue(ctx *Context, v Value) Value {
	return Value(C.JS_DupValue((*C.JSContext)(ctx), C.JSValue(v)))
}

func NewString(ctx *Context, str string) Value {
	cstr := C.CString(str)
	defer C.free(unsafe.Pointer(cstr))

	return Value(C.JS_NewString((*C.JSContext)(ctx), cstr))
}

func NewInt(ctx *Context, n int) Value {
	return Value(C.JS_NewInt32((*C.JSContext)(ctx), C.int32_t(n)))
}

func NewFloat(ctx *Context, n float64) Value {
	return Value(C.JS_NewFloat64((*C.JSContext)(ctx), C.double(n)))
}

func FreeValue(ctx *Context, v Value) {
	C.JS_FreeValue((*C.JSContext)(ctx), C.JSValue(v))
}

func GetException(ctx *Context) Value {
	return Value(C.JS_GetException((*C.JSContext)(ctx)))
}

func ParseJSON(ctx *Context, data string, filename string) Value {
	cdata := C.CString(data)
	defer C.free(unsafe.Pointer(cdata))

	cfilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	return Value(C.JS_ParseJSON((*C.JSContext)(ctx), cdata, csize(len(data)), cfilename))
}

func IsError(ctx *Context, v Value) bool {
	return C.JS_IsError((*C.JSContext)(ctx), C.JSValue(v)) != 0
}

func IsArray(ctx *Context, v Value) bool {
	return C.JS_IsArray((*C.JSContext)(ctx), C.JSValue(v)) != 0
}

func NewObject(ctx *Context) Value {
	return Value(C.JS_NewObject((*C.JSContext)(ctx)))
}

func NewObjectProto(ctx *Context, v Value) Value {
	return Value(C.JS_NewObjectProto((*C.JSContext)(ctx), C.JSValue(v)))
}

func Throw(ctx *Context, v Value) Value {
	return Value(C.JS_Throw((*C.JSContext)(ctx), C.JSValue(DupValue(ctx, v))))
}

func SetConstructorBit(ctx *Context, obj Value, val bool) {
	C.JS_SetConstructorBit((*C.JSContext)(ctx), C.JSValue(obj), C.int(booltoint(val)))
}

func NewAtom(ctx *Context, name string) Atom {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return Atom(C.JS_NewAtom((*C.JSContext)(ctx), cname))
}

func FreeAtom(ctx *Context, atom Atom) {
	C.JS_FreeAtom((*C.JSContext)(ctx), C.JSAtom(atom))
}

func AtomToValue(ctx *Context, atom Atom) Value {
	return Value(C.JS_AtomToValue((*C.JSContext)(ctx), C.JSAtom(atom)))
}

func AtomToString(ctx *Context, atom Atom) Value {
	return Value(C.JS_AtomToString((*C.JSContext)(ctx), C.JSAtom(atom)))
}

func SetConstructor(ctx *Context, funcObj Value, proto Value) {
	C.JS_SetConstructor((*C.JSContext)(ctx), C.JSValue(funcObj), C.JSValue(proto))
}

func IsJobPending(rt *Runtime) bool {
	return C.JS_IsJobPending((*C.JSRuntime)(rt)) != 0
}

func ExecutePendingJob(rt *Runtime) (*Context, int) {
	var ctx *Context
	return ctx, int(C.JS_ExecutePendingJob((*C.JSRuntime)(rt), (**C.JSContext)(unsafe.Pointer(&ctx))))
}

func IsFunction(ctx *Context, val Value) bool {
	return C.JS_IsFunction((*C.JSContext)(ctx), C.JSValue(val)) != 0
}

func EnqueueJob(ctx *Context, jobFunc, thisVal Value, args []Value) {
	args = append([]Value{jobFunc, thisVal}, args...)
	defer runtime.KeepAlive(args)

	jsargs := jsvalues(args)
	defer runtime.KeepAlive(jsargs)

	// JS_EnqueueJob returns -1 when malloc fails,
	// which should be dealt with by the custom malloc funcs
	C.JS_EnqueueJob((*C.JSContext)(ctx), (*C.JSJobFunc)(C.js_call_func_job), C.int(len(args)), argv(jsargs))
}

func RunGC(rt *Runtime) {
	C.JS_RunGC((*C.JSRuntime)(rt))
}

func WriteObject(ctx *Context, obj Value, flags WriteObjectFlag) []byte {
	var size csize
	ptr := C.JS_WriteObject((*C.JSContext)(ctx), &size, C.JSValue(obj), C.int(flags))
	defer C.js_free((*C.JSContext)(ctx), unsafe.Pointer(ptr))

	ret := make([]byte, int(size))
	copy(ret, *(*[]byte)(makeSliceHeader(unsafe.Pointer(ptr), int(size))))

	return ret
}

func ReadObject(ctx *Context, buf []byte, flags ReadObjectFlag) Value {
	return Value(C.JS_ReadObject((*C.JSContext)(ctx), (*C.uint8_t)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&buf)).Data)), csize(len(buf)), C.int(flags)))
}

func EvalFunction(ctx *Context, funObj Value) Value {
	return Value(C.JS_EvalFunction((*C.JSContext)(ctx), C.JSValue(DupValue(ctx, funObj))))
}
