package event

import (
	context "context"
	fmt "fmt"
	"log"

	"github.com/dop251/goja"
	"google.golang.org/grpc"
)

type powershellTask struct {
	callback   goja.Callable
	finished   bool
	scriptBody string
	response   string
}

func (current *powershellTask) HasFinished(vm *goja.Runtime) bool {
	if current.finished {
		current.callback(nil, vm.ToValue(current.response))
	}
	return current.finished
}

// NewPowershellTask cunstructor
func NewPowershellTask(call goja.Callable, scriptBody string) Task {
	fmt.Println("GRPC version ", grpc.Version)

	conn, err := grpc.Dial("localhost:5000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(fmt.Sprintf("%s.%s", "1", err))
	}

	client := NewGreeterClient(conn)
	res, err := client.SayHello(context.Background(), &HelloRequest{Name: "test Hello Request"})
	if err != nil {
		log.Fatal(fmt.Sprintf("%s.%s", "2", err))
	}

	task := &powershellTask{
		callback:   call,
		finished:   true,
		scriptBody: scriptBody,
		response:   res.GetMessage(),
	}

	return task
}
