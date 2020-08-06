package js

import "golang.org/x/sys/windows"

type threadID uint32

func currentThreadID() threadID {
	return threadID(windows.GetCurrentThreadId())
}
