package spsw

import (
	"errors"
)

type ExporterBackend interface {
	StartExporting(jobUUID string, fieldNames []string) error
	WriteItem(i *Item) error
	FinishExporting(jobUUID string) error
}

type AbstractExporterBackend struct {
	ExporterBackend
}

func (aeb *AbstractExporterBackend) StartExporting(jobUUID string, fieldNames []string) error {
	return errors.New("Not implemented")
}

func (aeb *AbstractExporterBackend) WriteItem(i *Item) error {
	return errors.New("Not implemented")
}

func (aeb *AbstractExporterBackend) FinishExporting(jobUUID string) error {
	return errors.New("Not implemented")
}
