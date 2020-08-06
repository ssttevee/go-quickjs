package internal

// #include "quickjs/quickjs.h"
import "C"
import "fmt"

var (
	classDefs = map[C.JSClassID]C.JSClassDef{}
)

func registerClassDefinition(id C.JSClassID, name string, def C.JSClassDef) {
	def.class_name = C.CString(name)
	classDefs[id] = def
}

func defineCustomClasses(rt *C.JSRuntime) *C.JSRuntime {
	for id, def := range classDefs {
		if C.JS_NewClass(rt, id, &def) != 0 {
			panic(fmt.Sprintf("failed to define custom class: %s", C.GoString(def.class_name)))
		}
	}

	return rt
}
