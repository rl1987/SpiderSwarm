package spiderswarm

type Worker struct {
	UUID string
	// XXX: perhaps this should be one channel per direction?
	TasksIn         chan *Task
	ItemsOut        chan *Item
	TaskPromisesOut chan *TaskPromise
	ErrorsOut       chan error
}

func (w *Worker) Run() error {
	for {
		// Get TaskPromise.
		// Run Tasks. One Task at a time, one Action at a time.
		// Send TaskPromise/Item/error to manager.

		break
	}

	return nil
}
