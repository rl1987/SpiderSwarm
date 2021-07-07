package main

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

type Vertex struct {
	Type        int
	Name        string
	XPath       string
	Regex       string
	CSSSelector string
	ItemName    string
	FieldName   string
}

type Edge struct {
	Name       string
	FromVertex string
	ToVertex   string
}

type SpiderTask struct {
	Name     string
	Vertices []Vertex
	Edges    []Edge
}

func main() {
	fmt.Println("SpiderSwarm")
	httpClient := &http.Client{}
	spew.Dump(httpClient)
}
