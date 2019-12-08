package event

import "github.com/dop251/goja"

// Task interface
type Task interface {
	HasFinished(vm *goja.Runtime) bool
}

// Loop struct
type Loop struct {
	workingTasks []Task
}

// NewLoop constructor
func NewLoop() *Loop {
	loop := &Loop{
		workingTasks: []Task{},
	}

	return loop
}

// Add new task to event loop
func (current *Loop) Add(task Task) {
	current.workingTasks = append(current.workingTasks, task)
}

// IsEmpty method returns true if there are no any waiting async tasks for this loop
func (current *Loop) IsEmpty(vm *goja.Runtime) bool {
	var notFinishedTasks []Task

	for i := range current.workingTasks {
		if current.workingTasks[i].HasFinished(vm) == false {
			notFinishedTasks = append(notFinishedTasks, current.workingTasks[i])
		}
	}

	current.workingTasks = notFinishedTasks

	return len(current.workingTasks) == 0
}
