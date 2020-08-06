package js

import "syscall"

type threadID int

func currentThreadID() threadID {
	return threadID(syscall.Gettid())
}
