package spiderswarm

type Worker struct {
	UUID string
	// XXX: perhaps this should be one channel per direction?
	TasksIn         chan *Task
	ItemsOut        chan *Item
	TaskPromisesOut chan *TaskPromise
	ErrorsOut       chan error
}
