package event

import (
	"time"

	"github.com/dop251/goja"
)

type oneTimeScheduler struct {
	runAfter time.Time
	callback goja.Callable
}

func (current *oneTimeScheduler) HasFinished() bool {
	if time.Now().After(current.runAfter) {
		current.callback(nil)
	}
	return time.Now().After(current.runAfter)
}

// NewOneTimeScheduler cunstructor
func NewOneTimeScheduler(call goja.Callable, afterMilliseconds int64) Task {
	scheduler := &oneTimeScheduler{
		runAfter: time.Now().Add(time.Millisecond * time.Duration(afterMilliseconds)),
		callback: call,
	}
	return scheduler
}
