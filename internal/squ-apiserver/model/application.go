package model

type Application struct {
	BaseModel
	Name        string
	Description string
	Type        string // compose and manifest
	Content     string
	Version     string
}

type ApplicationServer struct {
	BaseModel
	ServerID      uint
	ApplicationID uint
	Status        string // 应用在该服务器上的状态: running, stopped, error等
}
