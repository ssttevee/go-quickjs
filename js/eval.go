package js

import (
	"errors"
	"io/ioutil"
	"runtime"
	"strconv"

	"github.com/ssttevee/go-quickjs/internal"
)

func (rt *Runtime) nextVMName() string {
	rt.counter++
	return "VM:" + strconv.Itoa(rt.counter)
}

func (rt *Runtime) prepareNodeModuleRealm() (*Realm, *Value, error) {
	r, err := rt.NewRealm()
	if err != nil {
		return nil, nil, err
	}

	exports, err := r.NewObject()
	if err != nil {
		return nil, nil, err
	}

	globalObj, err := r.GlobalObject()
	if err != nil {
		return nil, nil, err
	}

	if ok, err := globalObj.Set("exports", exports); err != nil {
		return nil, nil, err
	} else if !ok {
		return nil, nil, errors.New("failed to define exports in global object")
	}

	return r, exports, nil
}

func (rt *Runtime) evalNodeModule(script, filename string, opts ...EvalOption) (*Value, error) {
	r, exports, err := rt.prepareNodeModuleRealm()
	if err != nil {
		return nil, err
	}

	if _, err := r.eval(script, filename, opts...); err != nil {
		return nil, err
	}

	return exports, nil
}

func (rt *Runtime) EvalNodeModule(script string, opts ...EvalOption) (*Value, error) {
	return rt.evalNodeModule(script, rt.nextVMName(), opts...)
}

func (rt *Runtime) EvalNodeModuleFile(file string, opts ...EvalOption) (*Value, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return rt.evalNodeModule(string(data), file, opts...)
}

func (rt *Runtime) EvalBinaryNodeModule(buf []byte) (*Value, error) {
	r, exports, err := rt.prepareNodeModuleRealm()
	if err != nil {
		return nil, err
	}

	if _, err := r.evalBinary(buf); err != nil {
		return nil, err
	}

	return exports, nil
}

func (r *Realm) eval(script, filename string, opts ...EvalOption) (*Value, error) {
	var config evalConfig
	for _, option := range opts {
		option(&config)
	}

	return r.createAndResolveValue(internal.Eval(r.context, script, filename, config.flags))
}

func (r *Realm) evalBinary(buf []byte) (*Value, error) {
	fn, err := r.createAndResolveValue(internal.ReadObject(r.context, buf, internal.ReadObjectBytecode))
	if err != nil {
		return nil, err
	}

	defer runtime.KeepAlive(fn)

	return r.createAndResolveValue(internal.EvalFunction(r.context, fn.value))
}

func (r *Realm) Eval(script string, opts ...EvalOption) (*Value, error) {
	return r.eval(script, r.runtime.nextVMName(), opts...)
}

func (r *Realm) EvalModule(script string, opts ...EvalOption) (*Value, error) {
	return r.eval(script, r.runtime.nextVMName(), append(opts, evalOptionModule)...)
}

func (r *Realm) EvalBinary(buf []byte) (*Value, error) {
	return r.evalBinary(buf)
}

func (r *Realm) evalFile(file string, opts ...EvalOption) (*Value, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return r.eval(string(data), file, opts...)
}

func (r *Realm) EvalFile(file string, opts ...EvalOption) (*Value, error) {
	return r.evalFile(file, opts...)
}

func (r *Realm) EvalModuleFile(file string, opts ...EvalOption) (*Value, error) {
	return r.evalFile(file, append(opts, evalOptionModule)...)
}
