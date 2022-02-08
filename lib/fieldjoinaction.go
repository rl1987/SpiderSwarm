package spsw

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

const FieldJoinActionOutputItem = "FieldJoinActionOutputItem"
const FieldJoinActionOutputMap = "FieldJoinActionOutputMap"

type FieldJoinAction struct {
	AbstractAction
	WorkflowName  string
	JobUUID       string
	TaskUUID      string
	ItemName      string
	RequireFields []string
}

// XXX: do we want to take all of these things as params? they seem to violate the
// abstraction here.
func NewFieldJoinAction(inputNames []string, jobUUID string, taskUUID string, itemName string) *FieldJoinAction {
	return &FieldJoinAction{
		AbstractAction: AbstractAction{
			AllowedInputNames:  inputNames,
			AllowedOutputNames: []string{FieldJoinActionOutputItem, FieldJoinActionOutputMap},
			Inputs:             map[string]*DataPipe{},
			Outputs:            map[string][]*DataPipe{},
			CanFail:            false,
			UUID:               uuid.New().String(),
		},
		JobUUID:  jobUUID,
		TaskUUID: taskUUID,
		ItemName: itemName,
	}
}

func NewFieldJoinActionFromTemplate(actionTempl *ActionTemplate) Action {
	var inputNames []string
	var itemName string

	inputNames = actionTempl.ConstructorParams["inputNames"].StringsValue
	itemName = actionTempl.ConstructorParams["itemName"].StringValue

	action := NewFieldJoinAction(inputNames, "", "", itemName)

	action.Name = actionTempl.Name

	if _, ok := actionTempl.ConstructorParams["requireFields"]; ok {
		action.RequireFields = actionTempl.ConstructorParams["requireFields"].StringsValue
	}

	return action
}

func (fja *FieldJoinAction) String() string {
	return fmt.Sprintf("<FieldJoinAction %s Name: %s WorkflowName: %s, ItemName: %s>", fja.UUID, fja.Name,
		fja.WorkflowName, fja.ItemName)
}

func (fja *FieldJoinAction) Run() error {
	if fja.Outputs[FieldJoinActionOutputItem] == nil && fja.Outputs[FieldJoinActionOutputMap] == nil {
		return errors.New("No output connected")
	}

	if len(fja.Inputs) == 0 {
		return errors.New("No inputs connected")
	}

	item := NewItem(fja.ItemName, fja.WorkflowName, fja.JobUUID, fja.TaskUUID)
	m := map[string]string{}

	for key, inDP := range fja.Inputs {
		if len(inDP.Queue) > 0 {
			value := inDP.Remove()
			item.SetField(key, value)
			s, ok := value.(string)
			if ok {
				m[key] = s
			}
		}
	}

	if fja.RequireFields != nil {
		for _, rf := range fja.RequireFields {
			if item.Fields[rf] == nil {
				return nil
			}
		}
	}

	if fja.Outputs[FieldJoinActionOutputItem] != nil {
		for _, outDP := range fja.Outputs[FieldJoinActionOutputItem] {
			outDP.AddItem(item)
		}
	}

	if fja.Outputs[FieldJoinActionOutputMap] != nil {
		for _, outDP := range fja.Outputs[FieldJoinActionOutputMap] {
			outDP.Add(m)
		}
	}

	return nil
}
