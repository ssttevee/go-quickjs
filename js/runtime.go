package js

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/ssttevee/go-quickjs/internal"
)

type Runtime struct {
	runtime *internal.Runtime

	defaultRealmOptions []RealmOption

	counter int

	taskQueue chan func() error

	mutex  sync.Mutex
	timers map[int]*time.Timer
}

func freeRuntime(rt *Runtime) {
	internal.FreeRuntime(rt.runtime)
}

func NewRuntime(defaultRealmOptions ...RealmOption) *Runtime {
	rt := &Runtime{
		runtime:             internal.NewRuntime(),
		defaultRealmOptions: defaultRealmOptions,
		timers:              map[int]*time.Timer{},
		taskQueue:           make(chan func() error, 512),
	}

	runtime.SetFinalizer(rt, freeRuntime)

	return rt
}

func (rt *Runtime) hasPendingTimer() bool {
	rt.mutex.Lock()
	defer rt.mutex.Unlock()

	return len(rt.timers) > 0
}

func (rt *Runtime) executePendingJob() (bool, error) {
	ctx, res := internal.ExecutePendingJob(rt.runtime)
	if res < 0 {
		return false, (&Realm{
			runtime: rt,
			context: ctx,
		}).getError()
	}

	return res != 0, nil
}

func (rt *Runtime) allocateTimerID() int {
	var id int
	for {
		id = int(rand.Int31())
		if _, ok := rt.timers[id]; !ok {
			return id
		}
	}
}

func (rt *Runtime) setTimeout(r *Realm, thisValue *Value, fn *Function, ms float64, args ...*Value) (int, error) {
	return rt.setTimer(r, fn, ms, args, func(id int) {
		rt.mutex.Lock()
		defer rt.mutex.Unlock()

		delete(rt.timers, id)
	})
}

func (rt *Runtime) setInterval(r *Realm, thisValue *Value, fn *Function, ms float64, args ...*Value) (int, error) {
	return rt.setTimer(r, fn, ms, args, func(id int) {
		rt.mutex.Lock()
		defer rt.mutex.Unlock()

		if timer, ok := rt.timers[id]; ok {
			timer.Reset(time.Duration(ms) * time.Millisecond)
		}
	})
}

func makeTaskFunc(r *Realm, fn, thisObject *Value, args interface{}) func() error {
	stack := debug.Stack()
	return func() error {
		var argValues []*Value
		switch args := args.(type) {
		case []*Value:
			argValues = args

		case []interface{}:
			var err error
			argValues, err = r.convertArgs(args)
			if err != nil {
				return err
			}

		default:
			panic(fmt.Sprintf("unexpected args type: %T", args))
		}

		// thisValue := internal.Undefined
		// if thisObject != nil {
		// 	thisValue = thisObject.value
		// }

		// defer runtime.KeepAlive(fn)
		// defer runtime.KeepAlive(thisObject)
		// defer runtime.KeepAlive(argValues)

		// internal.EnqueueJob(r.context, fn.value, thisValue, internalValues(argValues))

		if _, err := fn.CallValues(thisObject, argValues); err != nil {
			log.Println(string(stack))
			return err
		}

		return nil
	}
}

func (rt *Runtime) enqueueCall(r *Realm, fn, thisObject *Value, args []interface{}) {
	rt.taskQueue <- makeTaskFunc(r, fn, thisObject, args)
}

func (rt *Runtime) setTimer(r *Realm, fn *Function, ms float64, args []*Value, afterTask func(int)) (int, error) {
	rt.mutex.Lock()
	defer rt.mutex.Unlock()

	id := rt.allocateTimerID()
	task := makeTaskFunc(r, (*Value)(fn), nil, args)

	rt.timers[id] = time.AfterFunc(time.Duration(ms*float64(time.Millisecond)), func() {
		rt.taskQueue <- func() error {
			task()
			afterTask(id)

			return nil
		}
	})

	return id, nil
}

func (rt *Runtime) clearTimer(r *Realm, _ *Value, id int) {
	defer func() {
		rt.taskQueue <- func() error {
			// do nothing
			return nil
		}
	}()

	rt.mutex.Lock()
	defer rt.mutex.Unlock()

	defer delete(rt.timers, id)

	timer, ok := rt.timers[id]
	if ok {
		timer.Stop()
	}
}

func (rt *Runtime) HasAsyncTasks() bool {
	return internal.IsJobPending(rt.runtime) || rt.hasPendingTimer()
}

func (rt *Runtime) StartEventLoop(ctx context.Context, waitForever bool) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case task := <-rt.taskQueue:
			if err := task(); err != nil {
				return err
			}

		default:
		}

		ok, err := rt.executePendingJob()
		if err != nil {
			return err
		}

		runtime.GC()

		if !ok {
			if !waitForever && !rt.hasPendingTimer() {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()

			case task := <-rt.taskQueue:
				if err := task(); err != nil {
					return err
				}
			}
		}
	}
}

func (rt *Runtime) ParseJSON(data string, filename string) (*Value, error) {
	r, err := rt.NewRealm()
	if err != nil {
		return nil, err
	}

	return r.ParseJSON(data, filename)
}
