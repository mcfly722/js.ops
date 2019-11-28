package main

import (
	"fmt"
	"github.com/dop251/goja"
	"io/ioutil"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func main() {

	data, err := ioutil.ReadFile("sample.js")
	check(err)

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
	
	v, err := vm.RunString(string(data))
	if err != nil {
	    panic(err)
	}
	
	fmt.Println("Contents of js:", v)
}
