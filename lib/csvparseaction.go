package spsw

import (
  	"bytes"
	"fmt"
	"errors"
  	"encoding/csv"
	"io"

	"github.com/google/uuid"
)

const CSVParseActionInputCSVBytes = "CSVParseActionInputCSVBytes"
const CSVParseActionInputCSVStr = "CSVParseActionInputCSVStr"
const CSVParseActionOutputMap = "CSVParseActionOutputMap"

type CSVParseAction struct {
	AbstractAction
}

// TODO: optionally accept fieldNames argument and let header row be absent.
func NewCSVParseAction() *CSVParseAction {
	return &CSVParseAction{
		AbstractAction: AbstractAction{
			AllowedInputNames: []string{CSVParseActionInputCSVBytes, CSVParseActionInputCSVStr},
			AllowedOutputNames: []string{CSVParseActionOutputMap},
			Inputs: map[string]*DataPipe{},
			Outputs: map[string][]*DataPipe{},
			CanFail: false,
			UUID: uuid.New().String(),
		},
	}
}

func NewCSVParseActionFromTemplate(actionTempl *ActionTemplate) Action {
	action := NewCSVParseAction()

	action.Name = actionTempl.Name

	return action
}

func (cpa *CSVParseAction) String() string {
	return fmt.Sprintf("<CSVParseAction %s Name: %s>", cpa.UUID, cpa.Name)
}

func (cpa *CSVParseAction) Run() error {
	if cpa.Inputs[CSVParseActionInputCSVBytes] == nil && cpa.Inputs[CSVParseActionInputCSVStr] == nil {
		return errors.New("No input connected")
	}

	if cpa.Outputs[CSVParseActionOutputMap] == nil || len(cpa.Outputs[CSVParseActionOutputMap]) == 0 {
		return errors.New("No output connected")
	}

	var fieldNames []string
	var row []string
	var outputMap map[string][]string
	var csvBytes []byte
	
	if cpa.Inputs[CSVParseActionInputCSVBytes] != nil {
		csvBytes, _ = cpa.Inputs[CSVParseActionInputCSVBytes].Remove().([]byte)
	} else {
		csvStr, _ := cpa.Inputs[CSVParseActionInputCSVStr].Remove().(string)
		csvBytes = []byte(csvStr)
	}

	csvReader := csv.NewReader(bytes.NewReader(csvBytes))
	
	fieldNames, err := csvReader.Read()

	outputMap = map[string][]string{}

	for _, fieldName := range fieldNames {
		outputMap[fieldName] = []string{}
	}

	for {
		row, err = csvReader.Read()

		if row == nil {
			break
		}

		for i, field := range row {
			fieldName := fieldNames[i]

			outputMap[fieldName] = append(outputMap[fieldName], field)
		}
	}

	if err == io.EOF {
		err = nil
	}

	for _, outDP := range cpa.Outputs[CSVParseActionOutputMap] {
		outDP.Add(outputMap)
	}

	return err
}

