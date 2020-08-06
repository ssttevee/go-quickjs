package internal

// #include <stdlib.h>
// #include "quickjs/quickjs.h"
//
// extern void go_object_with_finalizer_class_finalizer(JSRuntime *rt, JSValue val);
import "C"
import (
	"math/rand"
	"sync"
	"unsafe"
)

var (
	objectWithFinalizerClassID = C.JS_NewClassID(new(C.JSClassID))

	savedFinalizersMutex sync.Mutex
	savedFinalizers      = map[int]func(){}
)

func lookupAndDeleteFinalizer(id int) func() {
	savedFinalizersMutex.Lock()
	defer savedFinalizersMutex.Unlock()

	// unpin from memory
	defer delete(savedFinalizers, id)

	return savedFinalizers[id]
}

func saveFinalizer(f func()) int {
	savedFinalizersMutex.Lock()
	defer savedFinalizersMutex.Unlock()

	// pin to memory
	var id int
	for {
		id = rand.Int()
		_, ok := savedFinalizers[id]
		if !ok {
			break
		}
	}

	savedFinalizers[id] = f

	return id
}

//export go_object_with_finalizer_class_finalizer
func go_object_with_finalizer_class_finalizer(runtime *C.JSRuntime, obj C.JSValue) {
	lookupAndDeleteFinalizer(int(uintptr(C.JS_GetOpaque(obj, objectWithFinalizerClassID))))()
}

func init() {
	registerClassDefinition(objectWithFinalizerClassID, "ObjectWithFinalizer", C.JSClassDef{
		finalizer: (*C.JSClassFinalizer)(C.go_object_with_finalizer_class_finalizer),
	})
}

func NewObjectWithFinalizer(ctx *Context, f func()) Value {
	obj := C.JS_NewObjectClass((*C.JSContext)(ctx), C.int(objectWithFinalizerClassID))
	C.JS_SetOpaque(obj, unsafe.Pointer(uintptr(saveFinalizer(f))))

	return Value(obj)
}
