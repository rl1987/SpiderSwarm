package main

type Workflow struct {
	Name    string
	Version string
	UUID    string
	Tasks   []Task // XXX: task templates?
}
