package js

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
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

	threadID threadID
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
		threadID:            currentThreadID(),
	}

	runtime.SetFinalizer(rt, freeRuntime)

	return rt
}

func (rt *Runtime) isSync() bool {
	return rt.threadID == currentThreadID()
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

type AsyncResult struct {
	Value *Value
	Error error
}

func makeTaskFunc(r *Realm, fn, thisObject *Value, args interface{}) (func() error, <-chan *AsyncResult) {
	result := make(chan *AsyncResult, 1)
	return func() error {
		var argValues []*Value
		switch args := args.(type) {
		case []*Value:
			argValues = args

		case []interface{}:
			var err error
			argValues, err = r.convertArgs(args)
			if err != nil {
				result <- &AsyncResult{Error: err}
				return nil
			}

		default:
			panic(fmt.Sprintf("unexpected args type: %T", args))
		}

		val, err := fn.CallValues(thisObject, argValues)
		result <- &AsyncResult{Value: val, Error: err}

		return nil
	}, result
}

func (rt *Runtime) enqueueCall(r *Realm, fn, thisObject *Value, args interface{}) <-chan *AsyncResult {
	task, result := makeTaskFunc(r, fn, thisObject, args)
	rt.enqueueTask(task)
	return result
}

func (rt *Runtime) enqueueTask(f func() error) {
	rt.taskQueue <- f
}

func (rt *Runtime) setTimer(r *Realm, fn *Function, ms float64, args []*Value, afterTask func(int)) (int, error) {
	rt.mutex.Lock()
	defer rt.mutex.Unlock()

	id := rt.allocateTimerID()
	task, resultChan := makeTaskFunc(r, (*Value)(fn), nil, args)

	rt.timers[id] = time.AfterFunc(time.Duration(ms*float64(time.Millisecond)), func() {
		rt.taskQueue <- func() error {
			if err := task(); err != nil {
				return err
			}

			result := <-resultChan
			if result.Error != nil {
				return result.Error
			}

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
