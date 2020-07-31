package js

import (
	"log"
)

type realmConfig struct {
	*Realm
}

type RealmOption func(r realmConfig) error

func makeLoggingFunc(logger *log.Logger, preArgs ...interface{}) interface{} {
	numPreArgs := len(preArgs)
	return func(r *Realm, _ *Value, args ...*Value) {
		values := make([]interface{}, len(args)+numPreArgs)
		for i, arg := range args {
			values[i+numPreArgs] = arg.String()
		}

		copy(values, preArgs)

		if logger == nil {
			log.Println(values...)
		} else {
			logger.Println(values...)
		}
	}
}

func AddIntrinsicConsole(logger ...*log.Logger) RealmOption {
	var l *log.Logger
	if len(logger) > 0 {
		l = logger[0]
	}

	return func(r realmConfig) error {
		consoleObj, err := r.NewObject()
		if err != nil {
			return err
		}

		if _, err := consoleObj.Set("log", makeLoggingFunc(l, "INFO:")); err != nil {
			return err
		}

		if _, err := consoleObj.Set("warn", makeLoggingFunc(l, "WARN:")); err != nil {
			return err
		}

		globalObj, err := r.GlobalObject()
		if err != nil {
			return err
		}

		if _, err := globalObj.Set("console", consoleObj); err != nil {
			return err
		}

		return nil
	}
}

func AddIntrinsicTimeout(r realmConfig) error {
	globalObj, err := r.GlobalObject()
	if err != nil {
		return err
	}

	if _, err := globalObj.Set("setTimeout", r.runtime.setTimeout); err != nil {
		return err
	}

	if _, err := globalObj.Set("setInterval", r.runtime.setInterval); err != nil {
		return err
	}

	if _, err := globalObj.Set("clearTimeout", r.runtime.clearTimer); err != nil {
		return err
	}

	if _, err := globalObj.Set("clearInterval", r.runtime.clearTimer); err != nil {
		return err
	}

	return nil
}
