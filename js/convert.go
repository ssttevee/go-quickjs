package js

import (
	"fmt"
	"reflect"

	"github.com/ssttevee/go-quickjs/internal"
)

type JSValuer interface {
	JSValue(r *Realm) (*Value, error)
}

type Literal string

func (l Literal) JSValue(r *Realm) (*Value, error) {
	return r.eval(string(l), r.runtime.nextVMName())
}

func (r *Realm) Convert(v interface{}) (*Value, error) {
	switch v := v.(type) {
	case *Value:
		return v, nil

	case JSValuer:
		return v.JSValue(r)

	case error:
		msg, err := r.Convert(v.Error())
		if err != nil {
			return nil, err
		}

		obj, err := r.createAndResolveValue(internal.NewError(r.context))
		if err != nil {
			return nil, err
		}

		if _, err := obj.Set("message", msg); err != nil {
			return nil, err
		}

		return obj, nil

	case string:
		return r.NewString(v)

	case int:
		return r.NewInt(v)

	case float64:
		return r.NewFloat(v)

	case bool:
		return r.NewBoolean(v)
	}

	if reflect.TypeOf(v).Kind() == reflect.Func {
		return r.NewFunction(v)
	}

	panic(fmt.Sprintf("type conversion not implemented for %T", v))
}

func (r *Realm) convertArgs(args []interface{}) ([]*Value, error) {
	convertedArgs := make([]*Value, len(args))
	for i, arg := range args {
		var err error
		convertedArgs[i], err = r.Convert(arg)
		if err != nil {
			return nil, err
		}
	}

	return convertedArgs, nil
}
