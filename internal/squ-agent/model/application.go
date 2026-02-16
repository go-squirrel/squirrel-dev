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
	Env         []map[string]string `gorm:"type:json;serializer:json"`
}

const (
	AppStatusStarting = "starting"
	AppStatusRunning  = "running"
	AppStatusStopped  = "stopped"
	AppStatusUndeploy = "undeploy"
	AppStatusPaused   = "paused"
	AppStatusFailed   = "Failed"
)
