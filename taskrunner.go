package spiderswarm

type TaskRunner struct {
	UUID            string
	TasksIn         chan *Task
	ItemsOut        chan *Item
	TaskPromisesOut chan *TaskPromise
	ErrorsOut       chan error
}
