package spiderswarm

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Exporter struct {
	UUID     string
	Backends []ExporterBackend
	ItemsIn  chan Item
}

func NewExporter() *Exporter {
	return &Exporter{
		UUID:     uuid.New().String(),
		Backends: []ExporterBackend{},
		ItemsIn:  make(chan Item),
	}
}

func (e *Exporter) Run() error {
	for item := range e.ItemsIn {
		// Receive items, pass them to exporter backend(s).
		log.Info(fmt.Printf("Exported %s got item %v", e.UUID, item))

		for _, backend := range e.Backends {
			backend.WriteItem(&item)
		}
	}

	return nil
}

func (e *Exporter) AddBackend(newBackend ExporterBackend) {
	e.Backends = append(e.Backends, newBackend)
}
