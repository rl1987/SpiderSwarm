package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	spsw "github.com/spiderswarm/spiderswarm/lib"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func runTestWorkflow() {
	workflow := &spsw.Workflow{}

	spew.Dump(workflow)

	fmt.Println("TODO")
}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func printUsage() {
	fmt.Println("Read the code for now")
}

func main() {
	initLogging()
	log.Info("Starting spiderswarm instance...")

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	singleNodeCmd := flag.NewFlagSet("singlenode", flag.ExitOnError)
	singleNodeWorkers := singleNodeCmd.Int("workers", 1, "Number of worker goroutines")

	switch os.Args[1] {
	case "singlenode":
		singleNodeCmd.Parse(os.Args[2:])
		log.Info(fmt.Sprintf("Number of worker goroutines: %d", *singleNodeWorkers))
	case "client":
		// TODO: client for REST API
		fmt.Println("client part not implemented yet")
	default:

		runTestWorkflow()
	}

}
