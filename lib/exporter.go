package spsw

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	UUID     string
	Backends []ExporterBackend
	ItemsIn  chan *Item
}

func NewExporter() *Exporter {
	return &Exporter{
		UUID:     uuid.New().String(),
		Backends: []ExporterBackend{},
		ItemsIn:  make(chan *Item),
	}
}

func (e *Exporter) Run() error {
	log.Info(fmt.Sprintf("Starting run loop for exporter %s", e.UUID))

	for item := range e.ItemsIn {
		// Receive items, pass them to exporter backend(s).
		log.Info(fmt.Sprintf("Exporter %s got item %v", e.UUID, item))

		for _, backend := range e.Backends {
			err := backend.WriteItem(item)
			if err != nil {
				log.Error(fmt.Sprintf("WriteItem failed with error: %v", err))
				spew.Dump(err)
			}
		}
	}

	return nil
}

func (e *Exporter) AddBackend(newBackend ExporterBackend) {
	e.Backends = append(e.Backends, newBackend)
}
