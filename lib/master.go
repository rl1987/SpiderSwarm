package spsw

type Master struct {
	UUID string
}

// TODO: implement REST API for workflow/job management
// * Create a job from workflow.
// * List running/finished/failed/all jobs.
// * Cancel a running job.
// * Get job with workflow and statistics.

func (m *Master) Run() error {
	//TODO: implement this

	for {
		// select across channels from managers
		// Tell managers about new tasks.
		// Keep state with statistics.
		// Make initial task promises when launching Workflow.
		// Get items from manager(s) and pass them to exporter.

		break
	}

	return nil
}
