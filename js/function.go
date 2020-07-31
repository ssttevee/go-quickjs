package js

import (
	"fmt"
	"log"
	"reflect"
	"runtime"
	"runtime/debug"

	"github.com/ssttevee/go-quickjs/internal"
)

type Function Value

func (f *Function) JSValue(r *Realm) (*Value, error) {
	return (*Value)(f), nil
}

func (f *Function) Call(thisObject *Value, args ...interface{}) (*Value, error) {
	return (*Value)(f).Call(thisObject, args...)
}

type TypedValue interface {
	FromValue(r *Realm, v *Value) (bool, error)
}

var (
	realmType      = reflect.TypeOf((*Realm)(nil))
	valueType      = reflect.TypeOf((*Value)(nil))
	functionType   = reflect.TypeOf((*Function)(nil))
	typedValueType = reflect.TypeOf((*TypedValue)(nil)).Elem()
	errorType      = reflect.TypeOf((*error)(nil)).Elem()
)

// type Function func(r *Realm, thisValue *Value, args []*Value) (*Value, error)

func jsToReflect(r *Realm, v *Value, t reflect.Type) (reflect.Value, error) {
	switch t {
	case functionType:
		return reflect.ValueOf((*Function)(v)), nil

	case valueType:
		return reflect.ValueOf(v), nil
	}

	if t.ConvertibleTo(typedValueType) {
		var argValue reflect.Value
		if t.Kind() == reflect.Ptr {
			argValue = reflect.New(t.Elem())
		} else {
			argValue = reflect.New(t).Elem()
		}

		ok, err := argValue.Interface().(TypedValue).FromValue(r, v)
		if err != nil {
			return reflect.Value{}, err
		}

		if !ok {
			return reflect.Value{}, &InvalidTypeError{Type: t}
		}

		return argValue, nil
	}

	switch t.Kind() {
	case reflect.String:
		if !v.IsString() {
			return reflect.Value{}, &InvalidTypeError{Type: t}
		}

		return reflect.ValueOf(v.ToString()), nil

	case reflect.Int:
		if !v.IsNumber() {
			return reflect.Value{}, &InvalidTypeError{Type: t}
		}

		return reflect.ValueOf(v.ToInt()), nil

	case reflect.Float64:
		return reflect.ValueOf(v.ToFloat()), nil

	case reflect.Slice:
		if !v.IsArray() {
			return reflect.Value{}, &InvalidTypeError{Type: t}
		}

		lengthValue, err := v.Get("length")
		if err != nil {
			return reflect.Value{}, err
		}

		n := lengthValue.ToInt()
		slice := reflect.MakeSlice(t, n, n)
		elemType := t.Elem()
		for i := 0; i < n; i++ {
			elem, err := v.Index(i)
			if err != nil {
				return reflect.Value{}, err
			}

			elemValue, err := jsToReflect(r, elem, elemType)
			if err != nil {
				return reflect.Value{}, err
			}

			runtime.KeepAlive(elem)

			slice.Index(i).Set(elemValue)
		}

		return slice, nil
	}

	panic(fmt.Sprintf("unexpected arg type: %s", t))
}

func prepareReflectCallArgs(r *Realm, thisValue internal.Value, args []internal.Value, argTypes []reflect.Type, variadic bool) ([]reflect.Value, error) {
	minArgs := len(argTypes)
	numCallArgs := len(argTypes)

	var variadicType reflect.Type
	if variadic {
		minArgs--
		numCallArgs = len(args)

		variadicType = argTypes[minArgs].Elem()
	}

	if len(args) < minArgs {
		return nil, &NotEnoughParametersError{
			Required: minArgs,
			Actual:   len(args),
		}
	}

	const argsOffset = 2
	callArgs := make([]reflect.Value, numCallArgs+argsOffset)

	for i, arg := range args {
		var argType reflect.Type
		if i < minArgs {
			argType = argTypes[i]
		} else if !variadic {
			break
		} else {
			argType = variadicType
		}

		var err error
		callArgs[i+argsOffset], err = jsToReflect(r, r.createValue(internal.DupValue(r.context, arg)), argType)
		if _, ok := err.(*InvalidTypeError); ok {
			return nil, &InvalidParameterTypeError{
				Index: i + 1,
				Type:  argType,
			}
		} else if err != nil {
			return nil, err
		}
	}

	callArgs[0] = reflect.ValueOf(r)
	callArgs[1] = reflect.ValueOf(r.createValue(internal.DupValue(r.context, thisValue)))

	return callArgs, nil
}

func (r *Realm) NewFunction(f interface{}) (*Value, error) {
	fValue := reflect.ValueOf(f)
	fType := fValue.Type()
	if fType.Kind() != reflect.Func {
		panic(fmt.Errorf("f must be a function, got %s", fType.Kind()))
	}

	const numPreArgs = 2

	numArgs := fType.NumIn()
	if numArgs < numPreArgs {
		panic(fmt.Errorf("f must have at least 2 args, got %d", numArgs))
	}

	if fType.In(0) != realmType {
		panic(fmt.Errorf("first arg of f must be *Realm, got %s", fType.In(0)))
	}

	if fType.In(1) != valueType {
		panic(fmt.Errorf("second arg of f must be *Value, got %s", fType.In(1)))
	}

	numOut := fType.NumOut()
	if numOut > 2 {
		panic(fmt.Errorf("f must have at most 2 return values, got %d", numOut))
	}

	if numOut > 0 && fType.Out(numOut-1) != errorType {
		panic(fmt.Errorf("last return value f must be error, got %s", fType.Out(numOut-1)))
	}

	argTypes := make([]reflect.Type, numArgs-numPreArgs)
	for i := numPreArgs; i < numArgs; i++ {
		argTypes[i-numPreArgs] = fType.In(i)
	}

	variadic := fType.IsVariadic()

	return r.createAndResolveValue(internal.NewFunction(r.context, func(ctx *internal.Context, thisValue internal.Value, args []internal.Value) internal.Value {
		callArgs, err := prepareReflectCallArgs(r, thisValue, args, argTypes, variadic)
		if err != nil {
			return r.throw(err)
		}

		out := fValue.Call(callArgs)
		if numOut > 0 {
			if err, ok := out[numOut-1].Interface().(error); ok && err != nil {
				errObj := Must(r.Convert(err))
				defer runtime.KeepAlive(errObj)

				return internal.Throw(ctx, errObj.value)
			}
		}

		if numOut < 2 || out[0].Kind() == reflect.Ptr && out[0].IsNil() {
			return internal.Undefined
		}

		result := Must(r.Convert(out[0].Interface()))
		defer runtime.KeepAlive(result)

		return internal.DupValue(ctx, result.value)
	}))
}
