package spiderswarm

type ExporterBackend interface {
	StartExporting(jobUUID string)
	WriteItem(i *Item) error
	FinishExporting(jobUUID string)
}

type AbstractExporterBackend struct {
	ExporterBackend
}
