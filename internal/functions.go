package internal

// #include <stdlib.h>
// #include "quickjs/quickjs.h"
//
// extern void go_function_class_finalizer(JSRuntime *rt, JSValue val);
// extern JSValue go_function_class_call(JSContext *ctx, JSValueConst func_obj, JSValueConst this_val, int argc, JSValueConst *argv, int flags);
import "C"
import (
	"math/rand"
	"sync"
	"unsafe"
)

var (
	goFunctionClassID = C.JS_NewClassID(new(C.JSClassID))

	savedFunctionsMutex sync.RWMutex
	savedFunctions      = map[int]Function{}
)

func defineGoFunctionClass(rt *C.JSRuntime) *C.JSRuntime {
	className := C.CString("GoFunction")
	defer C.free(unsafe.Pointer(className))

	if C.JS_NewClass(rt, goFunctionClassID, &C.JSClassDef{
		class_name: className,
		finalizer:  (*C.JSClassFinalizer)(C.go_function_class_finalizer),
		call:       (*C.JSClassCall)(C.go_function_class_call),
	}) != 0 {
		panic("failed to define go function class")
	}

	return rt
}

type Function func(ctx *Context, this Value, args []Value) Value

func funcID(obj C.JSValue) int {
	return int(uintptr(C.JS_GetOpaque(obj, goFunctionClassID)))
}

//export go_function_class_finalizer
func go_function_class_finalizer(runtime *C.JSRuntime, obj C.JSValue) {
	savedFunctionsMutex.Lock()
	defer savedFunctionsMutex.Unlock()

	// unpin from memory
	delete(savedFunctions, funcID(obj))
}

//export go_function_class_call
func go_function_class_call(ctx *C.JSContext, obj C.JSValue, thisValue C.JSValue, argc C.int, argv *C.JSValue, flags C.int) C.JSValue {
	args := make([]Value, int(argc))
	for i, arg := range *(*[]C.JSValue)(makeSliceHeader(unsafe.Pointer(argv), int(argc))) {
		args[i] = Value(arg)
	}

	return C.JSValue(lookupFunction(funcID(obj))((*Context)(ctx), Value(thisValue), args))
}

func lookupFunction(id int) Function {
	savedFunctionsMutex.RLock()
	defer savedFunctionsMutex.RUnlock()

	return savedFunctions[id]
}

func saveFunction(f Function) int {
	savedFunctionsMutex.Lock()
	defer savedFunctionsMutex.Unlock()

	// pin to memory
	var id int
	for {
		id = rand.Int()
		_, ok := savedFunctions[id]
		if !ok {
			break
		}
	}

	savedFunctions[id] = f

	return id
}

func NewFunction(ctx *Context, f Function) Value {
	funcObj := C.JS_NewObjectClass((*C.JSContext)(ctx), C.int(goFunctionClassID))
	C.JS_SetOpaque(funcObj, unsafe.Pointer(uintptr(saveFunction(f))))

	return Value(funcObj)
}
