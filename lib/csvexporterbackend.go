package spsw

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
)

// For now we're only supporting one item type per job. Later we want to implement
// multiple item types being exported during a single scraping job.

type CSVExporterBackend struct {
	AbstractExporterBackend

	OutputDirPath string

	csvWritersByJob  map[string]*csv.Writer
	fileHandlesByJob map[string]*os.File
	fieldNamesByJob  map[string][]string
}

func NewCSVExporterBackend(outputDirPath string) *CSVExporterBackend {
	if outputDirPath[len(outputDirPath)-1] == '/' {
		outputDirPath = outputDirPath[:len(outputDirPath)-1]
	}

	return &CSVExporterBackend{
		OutputDirPath:    outputDirPath,
		csvWritersByJob:  map[string]*csv.Writer{},
		fileHandlesByJob: map[string]*os.File{},
		fieldNamesByJob:  map[string][]string{},
	}
}

func (ceb *CSVExporterBackend) StartExporting(jobUUID string, fieldNames []string) error {
	// XXX: maybe include date/time into filename as well?
	csvFilePath := ceb.OutputDirPath + "/" + jobUUID + ".csv"

	csvFileHandle, err := os.OpenFile(csvFilePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		spew.Dump(err)
		return err
	}

	csvWriter := csv.NewWriter(csvFileHandle)
	err = csvWriter.Write(fieldNames)
	if err != nil {
		spew.Dump(err)
		csvFileHandle.Close()
		return err
	}

	ceb.csvWritersByJob[jobUUID] = csvWriter
	ceb.fileHandlesByJob[jobUUID] = csvFileHandle
	ceb.fieldNamesByJob[jobUUID] = fieldNames

	log.Info(fmt.Sprintf("Starting to export items to %s for job %s", csvFilePath, jobUUID))

	return nil
}

func (ceb *CSVExporterBackend) WriteItem(i *Item) error {
	jobUUID := i.JobUUID

	fieldNames := ceb.fieldNamesByJob[jobUUID]
	if fieldNames == nil {
		return errors.New("Field names not found in WriterItem for job " + jobUUID)
	}

	csvWriter := ceb.csvWritersByJob[jobUUID]
	if csvWriter == nil {
		return errors.New("CSV writer not found in WriteItem for job " + jobUUID)
	}

	row := []string{}

	for _, fieldName := range fieldNames {
		var rowStr string

		value := i.Fields[fieldName]

		if value.ValueType == ValueTypeString {
			rowStr = value.StringValue
		} else if value.ValueType == ValueTypeStrings {
			rowStr = "[" + strings.Join(value.StringsValue, ",") + "]"
		} else if value.ValueType == ValueTypeInt {
			rowStr = fmt.Sprintf("%d", value.IntValue)
		} else if value.ValueType == ValueTypeBool {
			if value.BoolValue {
				rowStr = "true"
			} else {
				rowStr = "false"
			}
		}

		row = append(row, rowStr)
	}

	err := csvWriter.WriteAll([][]string{row})
	if err != nil {
		return err
	}

	return nil
}

func (ceb *CSVExporterBackend) FinishExporting(jobUUID string) error {
	fileHandle := ceb.fileHandlesByJob[jobUUID]
	if fileHandle == nil {
		return errors.New("File handle not found in FinishExporting")
	}

	fileHandle.Sync()
	fileHandle.Close()

	ceb.csvWritersByJob[jobUUID] = nil
	ceb.fileHandlesByJob[jobUUID] = nil
	ceb.fieldNamesByJob[jobUUID] = nil

	log.Info(fmt.Sprintf("Finished exporting items for job %s", jobUUID))

	return nil
}
