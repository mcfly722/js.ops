package main

import (
	"fmt"
	"github.com/dop251/goja"
	"io/ioutil"
	"runtime"
	"sync"
	"./pkg"
	
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func process(scriptBody string, wg *sync.WaitGroup) {
	eventLoop := event.NewLoop()
	
	vm := goja.New()

	fconsole := func(msg string) {
		fmt.Println(msg)
	}

	fSetTimeout := func(call goja.Callable, delay goja.Value){
		fmt.Printf("Timeout configured for %v seconds\n", delay)
		call(nil,delay)
	}

	vm.Set("console", fconsole)
	vm.Set("setTimeout", fSetTimeout)
	
	v, err := vm.RunString(scriptBody)
	if err != nil {
	    panic(err)
	}
	
	fmt.Println("Contents of js:", v)
	
	// processing event loop till all events would be finished
	for ;!eventLoop.IsEmpty(); {
		runtime.Gosched();
	
	}
	
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	
	data, err := ioutil.ReadFile("sample.js")
	check(err)

	go process(string(data), &wg)

	wg.Wait()
    fmt.Print("done")
}
