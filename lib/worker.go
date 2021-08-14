package spiderswarm

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	UUID             string
	ScheduledTasksIn chan *ScheduledTask
	DataChunksOut    chan *DataChunk
}

func NewWorker() *Worker {
	return &Worker{
		UUID:             uuid.New().String(),
		ScheduledTasksIn: make(chan *ScheduledTask),
		DataChunksOut:    make(chan *DataChunk),
	}
}

func (w *Worker) executeTask(task *Task) error {
	err := task.Run()
	if err != nil {
		log.Error(fmt.Sprintf("Task %v failed with error: %v", task, err))
		// TODO: send error
		return err
	}

	for _, outDP := range task.Outputs {
		if len(outDP.Queue) == 0 {
			continue
		}

		x := outDP.Remove()

		if item, okItem := x.(*Item); okItem {
			for _, i := range item.Splay() {
				log.Info(fmt.Sprintf("Worker %s got item %v", w.UUID, i))
				chunk, _ := NewDataChunk(i)
				w.DataChunksOut <- chunk
			}
		}

		if promise, okPromise := x.(*TaskPromise); okPromise {
			for _, p := range promise.Splay() {
				log.Info(fmt.Sprintf("Worker %s enqueing promise %v", w.UUID, p))
				chunk, _ := NewDataChunk(p)
				w.DataChunksOut <- chunk
			}
		}
	}

	return nil
}

func (w *Worker) Run() error {
	for {
		scheduledTask := <-w.ScheduledTasksIn

		task := NewTaskFromScheduledTask(scheduledTask)
		log.Info(fmt.Sprintf("Worker %s running task %v", w.UUID, task))

		w.executeTask(task)
	}

	return nil
}
