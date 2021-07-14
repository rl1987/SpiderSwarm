package main

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

type DataPipe struct {
	Done  bool
	Queue []interface{}
}

type AbstractAction struct {
	Inputs     map[string]DataPipe
	Outputs    map[string]DataPipe
	CanFail    bool
	ExpectMany bool
}

type Action interface {
	Run() bool
}

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
		},
		BaseURL: baseURL,
		Method:  method,
	}
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
