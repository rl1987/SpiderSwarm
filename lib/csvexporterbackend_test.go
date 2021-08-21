package spiderswarm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCSVExporterBackend(t *testing.T) {
	outputDirPath := "/tmp/aaa/"

	backend := NewCSVExporterBackend(outputDirPath)

	assert.NotNil(t, backend)
	assert.Equal(t, "/tmp/aaa", backend.OutputDirPath)
}
