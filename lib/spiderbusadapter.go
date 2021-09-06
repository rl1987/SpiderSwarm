package spsw

import (
//"fmt"

//"github.com/google/uuid"
//log "github.com/sirupsen/logrus"
)

type SpiderBusAdapter struct {
	UUID              string
	Bus               *SpiderBus
	ScheduledTasksIn  chan *ScheduledTask
	ScheduledTasksOut chan *ScheduledTask
	TaskPromisesIn    chan *TaskPromise
	TaskPromisesOut   chan *TaskPromise
	ItemsIn           chan *Item
	ItemsOut          chan *Item
}
