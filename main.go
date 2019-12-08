package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"sync"

	event "./pkg"
	"github.com/dop251/goja"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func process(scriptBody string, wg *sync.WaitGroup) {

	wg.Add(1)

	go func() {
		eventLoop := event.NewLoop()

		vm := goja.New()

		fconsole := func(msg string) {
			fmt.Println(msg)
		}

		fSetTimeout := func(callback goja.Callable, delayMilliseconds goja.Value) {
			eventLoop.Add(event.NewOneTimeScheduler(callback, delayMilliseconds.ToInteger()))
		}

		fSetInterval := func(callback goja.Callable, intervalMilliseconds goja.Value) {
			eventLoop.Add(event.NewInfiniteTimeScheduler(callback, intervalMilliseconds.ToInteger()))
		}

		vm.Set("console", fconsole)
		vm.Set("setTimeout", fSetTimeout)
		vm.Set("setinterval", fSetInterval)

		v, err := vm.RunString(scriptBody)
		if err != nil {
			panic(err)
		}

		fmt.Println("Contents of js:", v)

		// processing event loop till all events would be finished
		for !eventLoop.IsEmpty(vm) {
			runtime.Gosched()
		}

		wg.Done()
	}()

}

func main() {
	var wg sync.WaitGroup

	data, err := ioutil.ReadFile("sample.js")
	check(err)

	process(string(data), &wg)

	wg.Wait()
	fmt.Println("done")
}
