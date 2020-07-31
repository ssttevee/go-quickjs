package js

import (
	"runtime"

	"github.com/ssttevee/go-quickjs/internal"
)

type propertyConfig struct {
	value, getter, setter *Value
	flags                 internal.PropertyFlag
}

type DefinePropertyOption func(*Realm, *propertyConfig) error

func DefinePropertyConfigurable(b bool) DefinePropertyOption {
	return func(r *Realm, c *propertyConfig) error {
		c.flags |= internal.PropertyFlagHasConfigurable

		if b {
			c.flags |= internal.PropertyFlagConfigurable
		} else {
			c.flags ^= internal.PropertyFlagConfigurable
		}

		return nil
	}
}

func DefinePropertyWritable(b bool) DefinePropertyOption {
	return func(r *Realm, c *propertyConfig) error {
		c.flags |= internal.PropertyFlagHasWritable

		if b {
			c.flags |= internal.PropertyFlagWritable
		} else {
			c.flags ^= internal.PropertyFlagWritable
		}

		return nil
	}
}

func DefinePropertyEnumerable(b bool) DefinePropertyOption {
	return func(r *Realm, c *propertyConfig) error {
		c.flags |= internal.PropertyFlagHasEnumerable

		if b {
			c.flags |= internal.PropertyFlagEnumerable
		} else {
			c.flags ^= internal.PropertyFlagEnumerable
		}

		return nil
	}
}

func DefinePropertyValue(v interface{}) DefinePropertyOption {
	return func(r *Realm, c *propertyConfig) error {
		c.flags |= internal.PropertyFlagHasValue

		var err error
		c.value, err = r.Convert(v)
		if err != nil {
			return err
		}

		return nil
	}
}

func DefinePropertyGetter(getter interface{}) DefinePropertyOption {
	return func(r *Realm, c *propertyConfig) error {
		c.flags |= internal.PropertyFlagHasGet

		var err error
		c.getter, err = r.Convert(getter)
		if err != nil {
			return err
		}

		return nil
	}
}

func DefinePropertySetter(setter interface{}) DefinePropertyOption {
	return func(r *Realm, c *propertyConfig) error {
		c.flags |= internal.PropertyFlagHasSet

		var err error
		c.setter, err = r.Convert(setter)
		if err != nil {
			return err
		}

		return nil
	}
}

func DefinePropertyThrow(c *propertyConfig) error {
	c.flags |= internal.PropertyFlagThrow
	return nil
}

func DefinePropertyNoExotic(c *propertyConfig) error {
	c.flags |= internal.PropertyFlagNoExotic
	return nil
}

func (v *Value) DefinePropertyAtom(prop *Atom, opts ...DefinePropertyOption) (bool, error) {
	defer runtime.KeepAlive(prop)

	var config propertyConfig
	defer runtime.KeepAlive(&config)

	for _, option := range opts {
		if err := option(v.realm, &config); err != nil {
			return false, err
		}
	}

	var (
		value  = internal.Undefined
		getter = internal.Undefined
		setter = internal.Undefined
	)

	if config.value != nil {
		value = config.value.value
	}

	if config.getter != nil {
		getter = config.getter.value
	}

	if config.setter != nil {
		setter = config.setter.value
	}

	result := internal.DefineProperty(v.realm.context, v.value, prop.atom, value, getter, setter, config.flags)
	if result == -1 {
		return false, v.realm.getError()
	}

	return result != 0, nil
}

func (v *Value) DefineProperty(prop string, opts ...DefinePropertyOption) (bool, error) {
	return v.DefinePropertyAtom(v.realm.NewStringAtom(prop), opts...)
}
