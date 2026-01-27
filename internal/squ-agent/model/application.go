package model

type Application struct {
	BaseModel
	Name        string
	Description string
	Type        string
	Status      string
	Content     string
	Version     string
	DeployID    uint64
}
