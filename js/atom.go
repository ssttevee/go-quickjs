package js

import (
	"runtime"

	"github.com/ssttevee/go-quickjs/internal"
)

type Atom struct {
	realm *Realm
	atom  internal.Atom
}

func freeAtom(a *Atom) {
	internal.FreeAtom(a.realm.context, a.atom)
	runtime.KeepAlive(a.realm)
}

func (r *Realm) createAtom(atom internal.Atom) *Atom {
	a := &Atom{
		realm: r,
		atom:  atom,
	}

	runtime.SetFinalizer(a, freeAtom)

	return a
}

func (r *Realm) NewStringAtom(s string) *Atom {
	return r.createAtom(internal.NewAtom(r.context, s))
}

func (a *Atom) ToString() (string, error) {
	value, err := a.ToStringValue()
	if err != nil {
		return "", err
	}

	return value.ToString(), nil
}

func (a *Atom) ToStringValue() (*Value, error) {
	return a.realm.createAndResolveValue(internal.AtomToString(a.realm.context, a.atom))
}

func (a *Atom) ToValue() (*Value, error) {
	return a.realm.createAndResolveValue(internal.AtomToValue(a.realm.context, a.atom))
}
