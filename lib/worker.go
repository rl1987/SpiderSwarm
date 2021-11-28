package spsw

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	UUID             string
	ScheduledTasksIn chan *ScheduledTask
	TaskPromisesOut  chan *TaskPromise
	TaskResultsOut   chan *TaskResult
	ItemsOut         chan *Item
	Done             chan interface{}
}

func NewWorker() *Worker {
	return &Worker{
		UUID:             uuid.New().String(),
		ScheduledTasksIn: make(chan *ScheduledTask),
		TaskPromisesOut:  make(chan *TaskPromise),
		TaskResultsOut:   make(chan *TaskResult),
		ItemsOut:         make(chan *Item),
		Done:             make(chan interface{}),
	}
}

func (w *Worker) String() string {
	return fmt.Sprintf("<Worker %s>", w.UUID)
}

func (w *Worker) executeTask(task *Task) error {
	err := task.Run()
	if err != nil {
		log.Error(fmt.Sprintf("Task %v failed with error: %v", task, err))

		taskResult := NewTaskResult(task.JobUUID, task.UUID, task.ScheduledTaskUUID, false, err)
		w.TaskResultsOut <- taskResult

		return err
	}

	taskResult := NewTaskResult(task.JobUUID, task.UUID, task.ScheduledTaskUUID, true, nil)

	nPromises := 0

	for outputName, outDP := range task.Outputs {
		if len(outDP.Queue) == 0 {
			continue
		}

		x := outDP.Remove()

		if item, okItem := x.(*Item); okItem {
			w.ItemsOut <- item

			taskResult.AddOutputItem(outputName, item)
		}

		if promise, okPromise := x.(*TaskPromise); okPromise {
			w.TaskPromisesOut <- promise
			nPromises++

			taskResult.AddOutputTaskPromise(outputName, promise)
		}
	}

	w.TaskResultsOut <- taskResult

	return nil
}

func (w *Worker) Run() error {
	log.Info(fmt.Sprintf("Starting runloop for worker %s", w.UUID))

	for {
		select {
		case scheduledTask := <-w.ScheduledTasksIn:
			if scheduledTask == nil {
				continue
			}

			log.Info(fmt.Printf("Worker %s got scheduled task %v", w.UUID, scheduledTask))

			task := NewTaskFromScheduledTask(scheduledTask)
			log.Info(fmt.Sprintf("Worker %s running task %v", w.UUID, task))

			w.executeTask(task)
		case <-w.Done:
			return nil
		}
	}

}
