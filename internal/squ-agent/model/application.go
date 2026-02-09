package model

type Application struct {
	BaseModel
	Name        string
	Description string
	Type        string
	OldStatus   string
	Status      string
	Content     string
	Version     string
	DeployID    uint64
}

const (
	AppStatusStarting = "starting"
	AppStatusRunning  = "running"
	AppStatusStopped  = "stopped"
	AppStatusUndeploy = "undeploy"
	AppStatusPaused   = "paused"
	AppStatusFailed   = "Failed"
)
