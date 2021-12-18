package spsw

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestExporterBackend struct {
	AbstractExporterBackend
	Items []*Item
}

func (teb *TestExporterBackend) WriteItem(item *Item) error {
	teb.Items = append(teb.Items, item)
	return nil
}

func TestExporterSimple(t *testing.T) {
	backend := &TestExporterBackend{
		Items: []*Item{},
	}

	exporter := NewExporter()

	exporter.AddBackend(backend)

	go exporter.Run()

	testItem := NewItem("testItem", "testWorkflow", "B927B203-5A25-44DA-AABB-3D2A41085B3F", "638585AD-8280-4990-8CDD-E8CFB6788D10")
	defer close(exporter.ItemsIn)

	exporter.ItemsIn <- testItem

	assert.Equal(t, 1, len(backend.Items))
	assert.Equal(t, testItem, backend.Items[0])
}
