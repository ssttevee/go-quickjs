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

func (r *Realm) NewStringAtom(s string) *Atom {
	a := &Atom{
		realm: r,
		atom:  internal.NewAtom(r.context, s),
	}

	runtime.SetFinalizer(a, freeAtom)

	return a
}
