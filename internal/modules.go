package internal

// #include "quickjs/quickjs.h"
//
// extern JSModuleDef *js_value_get_module_def(JSValue value);
//
// extern char *go_normalize_module(JSContext *ctx, char *module_base_name, char *module_name, void *opaque);
// extern JSModuleDef *go_load_module(JSContext *ctx, char *module_name, void *opaque);
import "C"
import (
	"math/rand"
	"unsafe"
)

type ModuleNormalizerFunc func(_ *Context, baseName, name string) string

type ModuleLoaderFunc func(_ *Context, name string) Value

var (
	normalizerFuncs = map[int]ModuleNormalizerFunc{}
	loaderFuncs     = map[int]ModuleLoaderFunc{}

	runtimeLoaderIDs = map[uintptr]int{}
)

//export go_normalize_module
func go_normalize_module(ctx *Context, baseName, name *C.char, opaque unsafe.Pointer) *C.char {
	normalizedName := normalizerFuncs[int(uintptr(opaque))](ctx, C.GoString(baseName), C.GoString(name))

	n := len(normalizedName)
	ret := C.js_malloc((*C.JSContext)(ctx), csize(n+1))
	*(*byte)(unsafe.Pointer(uintptr(ret) + uintptr(n))) = 0
	bytes := *(*[]byte)(makeSliceHeader(ret, n))
	copy(bytes, normalizedName)

	return (*C.char)(ret)
}

//export go_load_module
func go_load_module(ctx *Context, name *C.char, opaque unsafe.Pointer) *C.JSModuleDef {
	return C.js_value_get_module_def(C.JSValue(loaderFuncs[int(uintptr(opaque))](ctx, C.GoString(name))))
}

func SetModuleLoaderFunc(rt *Runtime, normalizeModule ModuleNormalizerFunc, loadModule ModuleLoaderFunc) {
	var id int
	for {
		id = rand.Int()
		if _, ok := normalizerFuncs[id]; ok {
			continue
		}

		if _, ok := loaderFuncs[id]; ok {
			continue
		}

		break
	}

	runtimeLoaderIDs[uintptr(unsafe.Pointer(rt))] = id

	var (
		normalizer *C.JSModuleNormalizeFunc
		loader     *C.JSModuleLoaderFunc
		opaque     unsafe.Pointer
	)

	if normalizeModule != nil {
		normalizerFuncs[id] = normalizeModule
		normalizer = (*C.JSModuleNormalizeFunc)(C.go_normalize_module)
	}

	if loadModule != nil {
		loaderFuncs[id] = loadModule
		loader = (*C.JSModuleLoaderFunc)(C.go_load_module)
	}

	if normalizer != nil && loader != nil {
		opaque = unsafe.Pointer(uintptr(id))
	}

	C.JS_SetModuleLoaderFunc((*C.JSRuntime)(rt), normalizer, loader, opaque)
}

func FreeModuleLoaderFunc(rt *Runtime) {
	id, ok := runtimeLoaderIDs[uintptr(unsafe.Pointer(rt))]
	if !ok {
		return
	}

	delete(normalizerFuncs, id)
	delete(loaderFuncs, id)
}
