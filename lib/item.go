package spsw

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type Item struct {
	UUID         string
	WorkflowName string
	JobUUID      string
	TaskUUID     string
	CreatedAt    time.Time
	Name         string

	Fields map[string]interface{}
}

func NewItem(name string, workflowName string, jobUUID string, taskUUID string) *Item {
	return &Item{
		UUID:         uuid.New().String(),
		WorkflowName: workflowName,
		JobUUID:      jobUUID,
		TaskUUID:     taskUUID,
		CreatedAt:    time.Now(),
		Fields:       map[string]interface{}{},
		Name:         name,
	}
}

func (i *Item) IsSplayable() bool {
	hasLists := false
	equalLen := true
	lastLen := -1

	for _, value := range i.Fields {
		rt := reflect.TypeOf(value) // XXX: this is bad for performance!

		if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
			hasLists = true

			if lastLen != -1 && lastLen != reflect.ValueOf(value).Len() {
				equalLen = false
				break
			}

			lastLen = reflect.ValueOf(value).Len()
		}
	}

	return hasLists && equalLen && lastLen != 0
}

func (i *Item) splayOff() *Item {
	newItem := &Item{
		UUID:         uuid.New().String(),
		WorkflowName: i.WorkflowName,
		JobUUID:      i.JobUUID,
		TaskUUID:     i.TaskUUID,
		CreatedAt:    time.Now(),
		Fields:       map[string]interface{}{},
		Name:         i.Name,
	}

	for key, value := range i.Fields {
		rt := reflect.TypeOf(value) // XXX: this is bad for performance!

		if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array {
			var x interface{}

			if s, ok1 := value.([]string); ok1 {
				x, i.Fields[key] = s[0], s[1:]
			} else {

				x, i.Fields[key] = value.([]interface{})[0], value.([]interface{})[1:]
			}
			newItem.Fields[key] = x
		} else {
			newItem.Fields[key] = value
		}
	}

	return newItem
}

func (i *Item) Splay() []*Item {
	if !i.IsSplayable() {
		return []*Item{i}
	}

	items := []*Item{}

	itemCopy := &Item{}

	copier.Copy(itemCopy, i)

	for {
		if !itemCopy.IsSplayable() {
			break
		}

		newItem := itemCopy.splayOff()

		items = append(items, newItem)
	}

	return items
}

func (i *Item) SetField(name string, value interface{}) {
	i.Fields[name] = interface{}(value)
}
