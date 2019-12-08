package event

import (
	"math"
	"time"

	"github.com/dop251/goja"
)

type oneTimeScheduler struct {
	runAfter time.Time
	callback goja.Callable
}

type infiniteTimeScheduler struct {
	startedAt            time.Time
	intervalMilliseconds int64
	firedCounter         int64
	callback             goja.Callable
}

func (current *oneTimeScheduler) HasFinished(vm *goja.Runtime) bool {
	if time.Now().After(current.runAfter) {
		current.callback(nil)
	}
	return time.Now().After(current.runAfter)
}

func (current *infiniteTimeScheduler) HasFinished(vm *goja.Runtime) bool {
	duration := time.Now().Sub(current.startedAt)

	if math.Floor((float64)(duration.Milliseconds())/(float64)(current.intervalMilliseconds)) > (float64)(current.firedCounter) {
		current.firedCounter = (int64)(math.Floor((float64)(duration.Milliseconds()) / (float64)(current.intervalMilliseconds)))
		current.callback(nil, vm.ToValue(current.firedCounter))
	}

	return false
}

// NewOneTimeScheduler cunstructor
func NewOneTimeScheduler(call goja.Callable, afterMilliseconds int64) Task {
	scheduler := &oneTimeScheduler{
		runAfter: time.Now().Add(time.Millisecond * time.Duration(afterMilliseconds)),
		callback: call,
	}
	return scheduler
}

// NewInfiniteTimeScheduler cunstructor
func NewInfiniteTimeScheduler(call goja.Callable, intervalMilliseconds int64) Task {

	if intervalMilliseconds < 1 {
		intervalMilliseconds = 1
	}

	scheduler := &infiniteTimeScheduler{
		startedAt:            time.Now(),
		intervalMilliseconds: intervalMilliseconds,
		firedCounter:         0,
		callback:             call,
	}
	return scheduler
}
