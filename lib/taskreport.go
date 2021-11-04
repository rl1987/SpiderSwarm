package spsw

type TaskReport struct {
	UUID      string
	JobUUID   string
	TaskUUID  string
	TaskName  string
	Succeeded bool
	Error     error
}
