package js

import (
	"io/ioutil"
	"runtime"

	"github.com/ssttevee/go-quickjs/internal"
)

func (rt *Runtime) Compile(script, filename string, opts ...EvalOption) ([]byte, error) {
	r, err := rt.NewRealm()
	if err != nil {
		return nil, err
	}

	v, err := r.createAndResolveValue(internal.Eval(r.context, script, filename, internal.EvalFlagCompileOnly))
	if err != nil {
		return nil, err
	}

	defer runtime.KeepAlive(v)

	return internal.WriteObject(r.context, v.value, internal.WriteObjectBytecode), nil
}

func (rt *Runtime) CompileModule(script, filename string, opts ...EvalOption) ([]byte, error) {
	return rt.Compile(script, filename, append(opts, evalOptionModule)...)
}

func (rt *Runtime) CompileFile(file string, opts ...EvalOption) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return rt.Compile(string(data), file, opts...)
}

func (rt *Runtime) CompileModuleFile(file string, opts ...EvalOption) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return rt.CompileModule(string(data), file, opts...)
}
