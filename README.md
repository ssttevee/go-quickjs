# go bindings for quickjs

Designed with go idioms in mind.

## Usage

```
import "github.com/ssttevee/go-quickjs/js"

...

r, _ := js.NewRuntime().NewRealm(js.AddIntrinsicConsole())
r.Eval(`console.log("hello world")`)
```

## Building libquickjs

### Windows (mingw-w64-x86_64)

```
pacman -Sy mingw-w64-x86_64-gcc mingw-w64-x86_64-pkg-config mingw-w64-x86_64-make mingw-w64-x86_64-dlfcn mingw-w64-x86_64-binutils
mingw32-make CONFIG_WIN64=y DESTDIR=../build install
```