package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

func initLogging() {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	initLogging()
	log.Info("Starting spiderswarm instance...")

	// https://apps.fcc.gov/cgb/form499/499a.cfm

	workflow := &Workflow{
		Name: "FCC_telecom",
		TaskTemplates: []TaskTemplate{
			TaskTemplate{
				TaskName: "ScrapeStates",
				Initial:  true,
			},
			TaskTemplate{
				TaskName: "ScrapeCompanyList",
				Initial:  false,
			},
			TaskTemplate{
				TaskName: "ScrapeCompanyPage",
				Initial:  false,
			},
		},
	}

	spew.Dump(workflow)
}
