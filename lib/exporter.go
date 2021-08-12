package spiderswarm

type Exporter struct {
	UUID     string
	Backends []ExporterBackend
}

func (e *Exporter) Run() error {
	for {
		// Receive items, pass them to exporter backend(s).

		break
	}

	return nil
}
