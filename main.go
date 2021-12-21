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
	fmt.Println("SpiderSwarm")
	fmt.Println("===========")
	fmt.Println("")
	fmt.Println("Run in single mode mode:")
	fmt.Println("  spiderswarm singlenode <backendAddr> <yamlFilePath>")
	fmt.Println("Run as worker with given number of worker goroutines:")
	fmt.Println("  spiderswarm worker <n> <backendAddr>")
	fmt.Println("Run as manager:")
	fmt.Println("  spiderswarm manager <backendAddr> <yamlFilePath>")
	fmt.Println("Run as exporter:")
	fmt.Println("  spiderswarm exporter <outputDir> <backendAddr>")
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
		if len(os.Args) != 4 {
			printUsage()
			os.Exit(0)
		}

		backendAddr := os.Args[2]
		yamlFilePath := os.Args[3]
		workflow := getWorkflow(yamlFilePath)
		success, err := workflow.Validate()
		if !success {
			fmt.Println(err)
			os.Exit(1)
		}

		runner.BackendAddr = backendAddr
		runner.RunSingleNode(4, ".", workflow)
		time.Sleep(1 * time.Second)
	case "worker":
		if len(os.Args) != 4 {
			printUsage()
			os.Exit(0)
		}

		n, _ := strconv.Atoi(os.Args[2])
		backendAddr := os.Args[3]
		runner.BackendAddr = backendAddr
		runner.RunWorkers(n)
		for {
			select {}
		}
	case "manager":
		if len(os.Args) != 4 {
			printUsage()
			os.Exit(0)
		}

		backendAddr := os.Args[2]
		yamlFilePath := os.Args[3]
		workflow := getWorkflow(yamlFilePath)

		success, err := workflow.Validate()
		if !success {
			fmt.Println(err)
			os.Exit(1)
		}

		runner.BackendAddr = backendAddr
		runner.RunManager(workflow)
		for {
			select {}
		}
	case "exporter":
		if len(os.Args) != 4 {
			printUsage()
			os.Exit(0)
		}

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
