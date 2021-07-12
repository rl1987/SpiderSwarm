package main

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"

	"github.com/rl1987/SpiderSwarm/workflow"
)

func main() {
	fmt.Println("SpiderSwarm")
	httpClient := &http.Client{}
	spew.Dump(httpClient)
}
