package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	spsw "github.com/spiderswarm/spiderswarm/lib"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v3"
)

func printUsage() {
	fmt.Println("Read the code for now")
}

func getWorkflow(filePath string) *spsw.Workflow {
	inF, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	decoder := yaml.NewDecoder(inF)

	workflow := &spsw.Workflow{}

	err = decoder.Decode(workflow)
	if err != nil {
		panic(err)
	}

	return workflow
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	runner := &spsw.Runner{}

	switch os.Args[1] {
	case "singlenode":
		backendAddr := os.Args[2]
		yamlFilePath := os.Args[3]
		workflow := getWorkflow(yamlFilePath)
		runner.BackendAddr = backendAddr
		runner.RunSingleNode(4, ".", workflow)
		time.Sleep(1 * time.Second)
	case "worker":
		n, _ := strconv.Atoi(os.Args[2])
		backendAddr := os.Args[3]
		runner.BackendAddr = backendAddr
		runner.RunWorkers(n)
		for {
			select {}
		}
	case "manager":
		backendAddr := os.Args[2]
		yamlFilePath := os.Args[3]
		workflow := getWorkflow(yamlFilePath)
		runner.BackendAddr = backendAddr
		runner.RunManager(workflow)
		for {
			select {}
		}
	case "exporter":
		outputDir := os.Args[2]
		backendAddr := os.Args[3]
		runner.BackendAddr = backendAddr
		runner.RunExporter(outputDir)
		for {
			select {}
		}
	case "client":
		// TODO: client for REST API
		log.Error("client part not implemented yet")
	default:
		printUsage()
	}
}
