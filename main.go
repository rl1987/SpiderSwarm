package main

import (
	"fmt"
	"net/http"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	fmt.Println("SpiderSwarm")
	httpClient := &http.Client{}
	spew.Dump(httpClient)
}
