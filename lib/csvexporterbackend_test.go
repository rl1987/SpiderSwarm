package spsw

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCSVExporterBackend(t *testing.T) {
	outputDirPath := "/tmp/aaa/"

	backend := NewCSVExporterBackend(outputDirPath)

	assert.NotNil(t, backend)
	assert.Equal(t, "/tmp/aaa", backend.OutputDirPath)
}

func TestCSVExporterBackendE2E(t *testing.T) {
	dir, err := ioutil.TempDir("", "spsw_test_")
	assert.Nil(t, err)

	defer os.RemoveAll(dir)

	jobUUID := "561E04F9-2A77-4B9F-ACB2-9A08AA3081CF"
	fieldNames := []string{"name", "phone"}

	backend := NewCSVExporterBackend(dir)

	_, err = backend.StartExporting(jobUUID, fieldNames)
	assert.Nil(t, err)

	csvFilePath := dir + "/" + jobUUID + ".csv"

	s, err := os.Stat(csvFilePath)
	assert.Nil(t, err)
	assert.Equal(t, 0755, int(s.Mode()))

	assert.Equal(t, 1, len(backend.csvWritersByJob))
	assert.Equal(t, 1, len(backend.fileHandlesByJob))
	assert.Equal(t, 1, len(backend.fieldNamesByJob))

	item1 := NewItem("person", "", jobUUID, "")
	item1.SetField("name", "Faust")
	item1.SetField("phone", "555-1212")

	item2 := NewItem("person", "", jobUUID, "")
	item2.SetField("name", "Mephistopheles")
	item2.SetField("phone", "666-0000")

	backend.WriteItem(item1)
	backend.WriteItem(item2)

	backend.FinishExporting(jobUUID)

	buf, err := ioutil.ReadFile(csvFilePath)
	csvStr := string(buf)

	expectedCsvStr := "name,phone\nFaust,555-1212\nMephistopheles,666-0000\n"

	assert.Equal(t, expectedCsvStr, csvStr)
}
