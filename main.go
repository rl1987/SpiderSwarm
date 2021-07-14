package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

type DataPipe struct {
	Done  bool
	Queue []interface{}
}

func NewDataPipe() *DataPipe {
	return &DataPipe{false, []interface{}{}}
}

type AbstractAction struct {
	Inputs             map[string]*DataPipe
	Outputs            map[string]*DataPipe
	CanFail            bool
	ExpectMany         bool
	AllowedInputNames  []string
	AllowedOutputNames []string
}

type Action interface {
	Run() bool
	AddInput(name string, dataPipe DataPipe)
	AddOutput(name string, dataPipe DataPipe)
}

const HTTPActionInputURLParams = "HTTPActionInputURLParams"
const HTTPActionInputHeaders = "HTTPActionInputHeaders"
const HTTPActionInputCookies = "HTTPActionInputCookies"

const HTTPActionOutputBody = "HTTPActionOutputBody"
const HTTPActionOutputHeaders = "HTTPActionOutputHeaders"
const HTTPActionOuputStatusCode = "HTTPActionOuputStatusCode"

type HTTPAction struct {
	AbstractAction
	BaseURL string
	Method  string
}

func NewHTTPAction(baseURL string, method string, canFail bool) *HTTPAction {
	return &HTTPAction{
		AbstractAction: AbstractAction{
			CanFail:    canFail,
			ExpectMany: false,
			AllowedInputNames: []string{
				HTTPActionInputURLParams,
				HTTPActionInputHeaders,
				HTTPActionInputCookies,
			},
			AllowedOutputNames: []string{
				HTTPActionOutputBody,
				HTTPActionOutputHeaders,
				HTTPActionOuputStatusCode,
			},
			Inputs:  map[string]*DataPipe{},
			Outputs: map[string]*DataPipe{},
		},
		BaseURL: baseURL,
		Method:  method,
	}
}

func (a *AbstractAction) AddInput(name string, dataPipe *DataPipe) error {
	for _, n := range a.AllowedInputNames {
		if n == name {
			a.Inputs[name] = dataPipe
			return nil
		}
	}

	return errors.New("input name not in AllowedInputNames")
}

func (a *AbstractAction) AddOutput(name string, dataPipe *DataPipe) error {
	for _, n := range a.AllowedOutputNames {
		if n == name {
			a.Outputs[name] = dataPipe
			return nil
		}
	}

	return errors.New("input name not in AllowedOutputNames")

}

type Task struct {
	Inputs  map[string]DataPipe
	Outputs map[string]DataPipe
}

type Workflow struct {
	Name    string
	Version string
	Tasks   []Task
}

func main() {
	fmt.Println("SpiderSwarm")
	httpClient := &http.Client{}
	spew.Dump(httpClient)

	w := Workflow{}
	spew.Dump(w)
}
