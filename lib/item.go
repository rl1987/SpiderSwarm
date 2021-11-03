package spsw

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	Fields map[string]*Value
}

func NewItem(name string, workflowName string, jobUUID string, taskUUID string) *Item {
	return &Item{
		UUID:         uuid.New().String(),
		WorkflowName: workflowName,
		JobUUID:      jobUUID,
		TaskUUID:     taskUUID,
		CreatedAt:    time.Now(),
		Fields:       map[string]*Value{},
		Name:         name,
	}
}

func NewItemFromJSON(raw []byte) *Item {
	item := &Item{}

	buffer := bytes.NewBuffer(raw)
	decoder := json.NewDecoder(buffer)

	err := decoder.Decode(item)
	if err != nil {
		return nil
	}

	return item
}

func (i *Item) Hash() []byte {
	h := sha256.New()

	h.Write([]byte(i.WorkflowName))
	h.Write([]byte(i.JobUUID))
	h.Write([]byte(i.Name))

	for key, value := range i.Fields {
		h.Write([]byte(key))
		h.Write(value.Hash())
	}

	return h.Sum(nil)
}

func (i *Item) String() string {
	return fmt.Sprintf("<Item %s WorkflowName: %s, JobUUID: %s, TaskUUID: %s, CreatedAt: %v, Name: %s, Fields: %v>",
		i.UUID, i.WorkflowName, i.JobUUID, i.TaskUUID, i.CreatedAt, i.Name, i.Fields)
}

func (i *Item) IsSplayable() bool {
	hasLists := false
	equalLen := true
	lastLen := -1

	for _, value := range i.Fields {
		if value.ValueType == ValueTypeStrings {
			hasLists = true

			if lastLen != -1 && lastLen != len(value.StringsValue) {
				equalLen = false
				break
			}

			lastLen = len(value.StringsValue)
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
		Fields:       map[string]*Value{},
		Name:         i.Name,
	}

	for key, value := range i.Fields {
		if value.ValueType == ValueTypeStrings {
			var x string
			x, value.StringsValue = value.StringsValue[0], value.StringsValue[1:]
			newItem.Fields[key] = NewValueFromString(x)
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
	if s, okStr := value.(string); okStr {
		i.Fields[name] = NewValueFromString(s)
	}

	if in, okInt := value.(int); okInt {
		i.Fields[name] = NewValueFromInt(in)
	}

	if str, okStr2 := value.([]string); okStr2 {
		i.Fields[name] = NewValueFromStrings(str)
	}

	if b, okBool := value.(bool); okBool {
		i.Fields[name] = NewValueFromBool(b)
	}
}

func (i *Item) EncodeToJSON() []byte {
	buffer := bytes.NewBuffer([]byte{})
	encoder := json.NewEncoder(buffer)

	encoder.Encode(i)

	bytes, _ := ioutil.ReadAll(buffer)

	return bytes
}

func (i *Item) FieldNames() []string {
	fieldNames := []string{}

	for key, _ := range i.Fields {
		fieldNames = append(fieldNames, key)
	}

	return fieldNames
}
