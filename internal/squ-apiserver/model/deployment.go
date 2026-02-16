package model

type Deployment struct {
	BaseModel
	ServerID      uint
	ApplicationID uint
	Status        string // 应用在该服务器上的状态: running, stopped, error等
	DeployID      uint64
	Content       string
	Env           []map[string]string `gorm:"type:json;serializer:json"`
}
