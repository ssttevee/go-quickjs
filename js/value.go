package js

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/ssttevee/go-quickjs/internal"
)

func NewNull() *Value {
	return &Value{value: internal.Null}
}

func NewUndefined() *Value {
	return &Value{value: internal.Undefined}
}

func NewFalse() *Value {
	return &Value{value: internal.False}
}

func NewTrue() *Value {
	return &Value{value: internal.True}
}

func NewBoolean(b bool) *Value {
	if b {
		return NewTrue()
	}

	return NewFalse()
}

type Value struct {
	realm *Realm
	value internal.Value
	// createStack []byte
}

func freeValue(v *Value) {
	internal.FreeValue(v.realm.context, v.value)
	runtime.KeepAlive(v.realm)
}

func (r *Realm) createValue(value internal.Value) *Value {
	v := &Value{
		realm: r,
		value: value,
		// createStack: debug.Stack(),
	}

	runtime.SetFinalizer(v, freeValue)

	return v
}

func (r *Realm) resolveValue(v *Value) (*Value, error) {
	if v.Tag() == TagException {
		return nil, r.getError()
	}

	return v, nil
}

func (r *Realm) createAndResolveValue(v internal.Value) (*Value, error) {
	return r.resolveValue(r.createValue(v))
}

func (v *Value) Tag() Tag {
	return Tag(v.value.Tag())
}

func (v *Value) IsString() bool {
	return v.value.Tag() == internal.TagString
}

func (v *Value) IsObject() bool {
	return v.value.Tag() == internal.TagObject
}

func (v *Value) IsArray() bool {
	defer runtime.KeepAlive(v)

	return internal.IsArray(v.realm.context, v.value)
}

func (v *Value) IsFunction() bool {
	defer runtime.KeepAlive(v)

	return internal.IsFunction(v.realm.context, v.value)
}

func (v *Value) IsInt() bool {
	return v.value.Tag() == internal.TagInt
}

func (v *Value) IsFloat() bool {
	return v.value.Tag() == internal.TagFloat64
}

func (v *Value) IsNumber() bool {
	return v.IsInt() || v.IsFloat()
}

func (v *Value) Interface() interface{} {
	switch v.Tag() {
	case TagString:
		return v.ToString()

	case TagInt:
		return v.ToInt()

	case TagBool:
		return v.ToBool()

	case TagNull, TagUndefined:
		return nil

	case TagFloat64:
		return v.ToFloat()
	}

	panic("unexpected value type " + Tag(v.Tag()).String())
}

func (v *Value) String() string {
	defer runtime.KeepAlive(v)

	switch v.Tag() {
	case TagNull:
		return "null"

	case TagUndefined:
		return "undefined"

	case TagObject:
		if v.IsArray() {
			n := Must(v.Get("length")).ToInt()

			elems := make([]string, n)
			for i := 0; i < n; i++ {
				elems[i] = Must(v.Index(i)).String()
			}

			return "[" + strings.Join(elems, " ") + "]"
		}

		if toStringFunc := Must(v.Get("toString")); toStringFunc.IsFunction() {
			return Must(toStringFunc.Call(v)).ToString()
		}

		className := Must(Must(v.Get("constructor")).Get("name")).ToString()
		if className == "" {
			className = "Object"
		}

		return "[object " + className + "]"
	}

	switch v := v.Interface().(type) {
	case string:
		return v

	case int:
		return strconv.Itoa(v)

	case bool:
		return strconv.FormatBool(v)

	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64)

	default:
		panic(fmt.Sprintf("unexpected value type %T", v))
	}
}

func (v *Value) ToString() string {
	if v.value.Tag() != internal.TagString {
		panic("value is not string")
	}

	return internal.ToString(v.realm.context, v.value)
}

func (v *Value) ToInt() int {
	defer runtime.KeepAlive(v)

	switch v.value.Tag() {
	case internal.TagInt:
		return internal.ToInt(v.value)

	case internal.TagFloat64:
		return int(internal.ToFloat(v.value))
	}

	panic("value is not number")
}

func (v *Value) ToFloat() float64 {
	defer runtime.KeepAlive(v)

	switch v.value.Tag() {
	case internal.TagInt:
		return float64(internal.ToInt(v.value))

	case internal.TagFloat64:
		return internal.ToFloat(v.value)
	}

	panic("value is not number")
}

func (v *Value) ToBool() bool {
	defer runtime.KeepAlive(v)

	switch v.value.Tag() {
	case internal.TagBool:
		return internal.ToBool(v.value)
	}

	panic("value is not boolean")
}

func (v *Value) IsTruthy() (bool, error) {
	defer runtime.KeepAlive(v)

	result := internal.IsTruthy(v.realm.context, v.value)
	if result == -1 {
		return false, v.realm.getError()
	}

	return result != 0, nil
}

func (v *Value) Get(property string) (*Value, error) {
	defer runtime.KeepAlive(v)

	return v.realm.createAndResolveValue(internal.GetPropertyStr(v.realm.context, v.value, property))
}

func (v *Value) Index(index int) (*Value, error) {
	defer runtime.KeepAlive(v)

	return v.realm.createAndResolveValue(internal.GetPropertyInt(v.realm.context, v.value, index))
}

func (v *Value) CallValues(thisObject *Value, args []*Value) (*Value, error) {
	thisValue := internal.Undefined
	if thisObject != nil {
		thisValue = thisObject.value
	}

	defer runtime.KeepAlive(v)
	defer runtime.KeepAlive(thisObject)
	defer runtime.KeepAlive(args)

	return v.realm.createAndResolveValue(internal.Call(v.realm.context, v.value, thisValue, internalValues(args)))
}

func (v *Value) Call(thisObject *Value, args ...interface{}) (*Value, error) {
	convertedArgs, err := v.realm.convertArgs(args)
	if err != nil {
		return nil, err
	}

	return v.CallValues(thisObject, convertedArgs)
}

func (v *Value) CallValuesAsync(thisObject *Value, args []*Value) <-chan *AsyncResult {
	return v.realm.runtime.enqueueCall(v.realm, v, thisObject, args)
}

func (v *Value) CallAsync(thisObject *Value, args ...interface{}) <-chan *AsyncResult {
	return v.realm.runtime.enqueueCall(v.realm, v, thisObject, args)
}

func (v *Value) InvokeValues(name string, args []*Value) (*Value, error) {
	defer runtime.KeepAlive(v)
	defer runtime.KeepAlive(args)

	if v.realm.runtime.isSync() {
		return v.realm.createAndResolveValue(internal.InvokeStr(v.realm.context, v.value, name, internalValues(args)))
	}

	result := <-v.InvokeValuesAsync(name, args)
	return result.Value, result.Error
}

func (v *Value) Invoke(name string, args ...interface{}) (*Value, error) {
	if !v.realm.runtime.isSync() {
		result := <-v.InvokeAsync(name, args...)
		return result.Value, result.Error
	}

	convertedArgs, err := v.realm.convertArgs(args)
	if err != nil {
		return nil, err
	}

	return v.InvokeValues(name, convertedArgs)
}

func (v *Value) InvokeValuesAsync(name string, args []*Value) <-chan *AsyncResult {
	funcValue, err := v.Get(name)
	if err != nil {
		result := make(chan *AsyncResult, 1)
		result <- &AsyncResult{Error: err}
		return result
	}

	return funcValue.CallValuesAsync(v, args)
}

func (v *Value) InvokeAsync(name string, args ...interface{}) <-chan *AsyncResult {
	funcValue, err := v.Get(name)
	if err != nil {
		result := make(chan *AsyncResult, 1)
		result <- &AsyncResult{Error: err}
		return result
	}

	return funcValue.CallAsync(v, args...)
}

func (v *Value) Construct(args ...interface{}) (*Value, error) {
	convertedArgs, err := v.realm.convertArgs(args)
	if err != nil {
		return nil, err
	}

	defer runtime.KeepAlive(v)
	defer runtime.KeepAlive(convertedArgs)

	return v.realm.createAndResolveValue(internal.CallConstructor(v.realm.context, v.value, internalValues(convertedArgs)))
}

func (v *Value) Set(prop string, val interface{}) (bool, error) {
	convertedValue, err := v.realm.Convert(val)
	if err != nil {
		return false, err
	}

	defer runtime.KeepAlive(v)
	defer runtime.KeepAlive(convertedValue)

	result := internal.SetPropertyStr(v.realm.context, v.value, prop, convertedValue.value)
	if result == -1 {
		return false, v.realm.getError()
	}

	return result != 0, nil
}

// WriteTo writes a pre-compiled script to the writer
func (v *Value) Bytes() []byte {
	defer runtime.KeepAlive(v)

	return internal.WriteObject(v.realm.context, v.value, internal.WriteObjectBytecode)
}
