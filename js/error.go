package js

import (
	"fmt"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/ssttevee/go-quickjs/internal"

	"github.com/dustin/go-humanize/english"
)

type Error Value

func (e *Error) JSValue(r *Realm) (*Value, error) {
	return (*Value)(e), nil
}

func (r *Realm) getError() error {
	v := internal.GetException(r.context)
	if v.Tag() == internal.TagNull {
		return nil
	}

	debug.PrintStack()

	return (*Error)(r.createValue(v))
}

func (e *Error) Error() string {
	if internal.IsError(e.realm.context, e.value) {
		return Must((*Value)(e).Get("name")).String() + ": " + Must((*Value)(e).Get("message")).String()
	}

	tag := (*Value)(e).Tag()
	switch tag {
	case TagUndefined, TagNull:
		return tag.String()
	}

	return Must((*Value)(e).Invoke("toString")).String()
}

func (e *Error) Stack() string {
	stack := e.Error()
	if internal.IsError(e.realm.context, e.value) {
		stack += "\n" + Must((*Value)(e).Get("stack")).String()
	}

	return strings.TrimSpace(stack)
}

type SyntaxError string

func NewSyntaxError(format string, v ...interface{}) error {
	return SyntaxError(fmt.Sprintf(format, v...))
}

func (e SyntaxError) Error() string {
	return string(e)
}

type TypeError string

func NewTypeError(format string, v ...interface{}) error {
	return TypeError(fmt.Sprintf(format, v...))
}

func (e TypeError) Error() string {
	return string(e)
}

type ReferenceError string

func NewReferenceError(format string, v ...interface{}) error {
	return ReferenceError(fmt.Sprintf(format, v...))
}

func (e ReferenceError) Error() string {
	return string(e)
}

type RangeError string

func NewRangeError(format string, v ...interface{}) error {
	return RangeError(fmt.Sprintf(format, v...))
}

func (e RangeError) Error() string {
	return string(e)
}

type InternalError string

func NewInternalError(format string, v ...interface{}) error {
	return InternalError(fmt.Sprintf(format, v...))
}

func (e InternalError) Error() string {
	return string(e)
}

type InvalidTypeError struct {
	Type reflect.Type
}

func (e *InvalidTypeError) isTypeError() {}

func (e *InvalidTypeError) Error() string {
	return fmt.Sprintf("value is not of type '%s'", e.Type)
}

type InvalidParameterTypeError struct {
	Index int
	Type  reflect.Type
}

func (e *InvalidParameterTypeError) isTypeError() {}

func (e *InvalidParameterTypeError) Error() string {
	return fmt.Sprintf("parameter %d is not of type '%s'", e.Index, e.Type)
}

type NotEnoughParametersError struct {
	Required, Actual int
}

func (e *NotEnoughParametersError) isTypeError() {}

func (e *NotEnoughParametersError) Error() string {
	return fmt.Sprintf("%s required, but only %d present.", english.Plural(e.Required, "argument", "arguments"), e.Actual)
}
