package js

import (
	"errors"
	"runtime"

	"github.com/ssttevee/go-quickjs/internal"
)

func Must(v *Value, err error) *Value {
	if err != nil {
		panic(err)
	}

	return v
}

func internalValues(values []*Value) []internal.Value {
	ret := make([]internal.Value, len(values))
	for i, value := range values {
		ret[i] = value.value
	}

	return ret
}

func valuesToInterfaces(values []*Value) []interface{} {
	ret := make([]interface{}, len(values))
	for i, value := range values {
		ret[i] = value.value
	}

	return ret
}

func ObjectAssign(dst *Value, src *Value) (*Value, error) {
	if !dst.IsObject() || !src.IsObject() {
		return nil, errors.New("dst and src must be objects")
	}

	defer runtime.KeepAlive(src)
	defer runtime.KeepAlive(dst)

	for _, prop := range internal.GetOwnPropertyNames(src.realm.context, src.value, 0) {
		atom := prop.Atom()
		internal.SetProperty(dst.realm.context, dst.value, atom, internal.GetProperty(src.realm.context, src.value, atom))
	}

	return dst, nil
}
