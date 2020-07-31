package js

import (
	"errors"
	"strconv"

	"github.com/ssttevee/go-quickjs/internal"
)

type Tag internal.Tag

const (
	TagBigDecimal    = Tag(internal.TagBigDecimal)
	TagBigInt        = Tag(internal.TagBigInt)
	TagBigFloat      = Tag(internal.TagBigFloat)
	TagSymbol        = Tag(internal.TagSymbol)
	TagString        = Tag(internal.TagString)
	TagObject        = Tag(internal.TagObject)
	TagInt           = Tag(internal.TagInt)
	TagBool          = Tag(internal.TagBool)
	TagNull          = Tag(internal.TagNull)
	TagUndefined     = Tag(internal.TagUndefined)
	TagUninitialized = Tag(internal.TagUninitialized)
	TagCatchOffset   = Tag(internal.TagCatchOffset)
	TagException     = Tag(internal.TagException)
	TagFloat64       = Tag(internal.TagFloat64)
)

func (t Tag) String() string {
	switch internal.Tag(t) {
	case internal.TagBigDecimal:
		return "bigdecimal"
	case internal.TagBigInt:
		return "bigint"
	case internal.TagBigFloat:
		return "bigfloat"
	case internal.TagSymbol:
		return "symbol"
	case internal.TagString:
		return "string"
	case internal.TagModule:
		return "module"
	case internal.TagFunctionBytecode:
		return "functionbytecode"
	case internal.TagObject:
		return "object"
	case internal.TagInt:
		return "int"
	case internal.TagBool:
		return "boolean"
	case internal.TagNull:
		return "null"
	case internal.TagUndefined:
		return "undefined"
	case internal.TagUninitialized:
		return "uninitialized"
	case internal.TagCatchOffset:
		return "offset"
	case internal.TagException:
		return "exception"
	case internal.TagFloat64:
		return "float64"
	}

	panic(errors.New("unexpected tag " + strconv.Itoa(int(t))))
}
