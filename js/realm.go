package js

import (
	"runtime"

	"github.com/ssttevee/go-quickjs/internal"
)

// Realm is a distinct global environment.
//
// See https://github.com/tc39/proposal-realms.
type Realm struct {
	runtime *Runtime
	context *internal.Context
}

func freeRealm(r *Realm) {
	internal.FreeContext(r.context)
	runtime.KeepAlive(r.runtime)
}

func (rt *Runtime) NewRealm(opts ...RealmOption) (*Realm, error) {
	r := &Realm{
		runtime: rt,
		context: internal.NewContext(rt.runtime),
	}

	runtime.SetFinalizer(r, freeRealm)

	for _, option := range rt.defaultRealmOptions {
		if err := option(realmConfig{r}); err != nil {
			return nil, err
		}
	}

	for _, option := range opts {
		if err := option(realmConfig{r}); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (r *Realm) NewObjectProto(proto *Value) (*Value, error) {
	return r.createAndResolveValue(internal.NewObjectProto(r.context, proto.value))
}

func (r *Realm) NewObject() (*Value, error) {
	return r.createAndResolveValue(internal.NewObject(r.context))
}

func (r *Realm) NewObjectWithFinalizer(f func()) (*Value, error) {
	return r.createAndResolveValue(internal.NewObjectWithFinalizer(r.context, f))
}

func (r *Realm) NewString(s string) (*Value, error) {
	return r.createAndResolveValue(internal.NewString(r.context, s))
}

func (r *Realm) NewInt(n int) (*Value, error) {
	return r.createAndResolveValue(internal.NewInt(r.context, n))
}

func (r *Realm) NewFloat(n float64) (*Value, error) {
	return r.createAndResolveValue(internal.NewFloat(r.context, n))
}

func (r *Realm) NewBoolean(b bool) (*Value, error) {
	return NewBoolean(b), nil
}

func (r *Realm) NewArrayBuffer(data []byte) (*Value, error) {
	return r.createAndResolveValue(internal.NewArrayBuffer(r.context, data))
}

func (r *Realm) SetConstructor(funcObj, proto *Value) {
	internal.SetConstructor(r.context, funcObj.value, proto.value)
}

func (r *Realm) throw(err error) internal.Value {
	switch err := err.(type) {
	case SyntaxError:
		return internal.ThrowSyntaxError(r.context, err.Error())
	case TypeError, interface {
		error
		isTypeError()
	}:
		return internal.ThrowTypeError(r.context, err.Error())
	case ReferenceError:
		return internal.ThrowReferenceError(r.context, err.Error())
	case RangeError:
		return internal.ThrowRangeError(r.context, err.Error())
	case InternalError:
		return internal.ThrowInternalError(r.context, err.Error())
	}

	errObj := Must(r.Convert(err))
	defer runtime.KeepAlive(errObj)

	return internal.Throw(r.context, errObj.value)
}

func (r *Realm) SetConstructorBit(obj *Value, val bool) {
	internal.SetConstructorBit(r.context, obj.value, val)
}

func (r *Realm) GlobalObject() (*Value, error) {
	return r.createAndResolveValue(internal.GetGlobalObject(r.context))
}

func (r *Realm) ParseJSON(data string, filename string) (*Value, error) {
	return r.createAndResolveValue(internal.ParseJSON(r.context, data, filename))
}

func (r *Realm) LoadValue(buf []byte) (*Value, error) {
	return r.createAndResolveValue(internal.ReadObject(r.context, buf, internal.ReadObjectBytecode))
}
