package event

type Task interface {
	HasFinished() bool
}

type Loop struct {
	WorkingTasks []Task
}

func NewLoop() *Loop {
	loop := &Loop{
		WorkingTasks: []Task{},
	}
	
	return loop;
}

func (current *Loop) IsEmpty() bool {
	var notFinishedTasks []Task

	for i := range current.WorkingTasks {
		if current.WorkingTasks[i].HasFinished() == false {
			notFinishedTasks = append(notFinishedTasks, current.WorkingTasks[i])
		}
	}

	current.WorkingTasks = notFinishedTasks

	return len(current.WorkingTasks) == 0
}
