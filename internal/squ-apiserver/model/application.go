package model

type Application struct {
	BaseModel
	Name        string
	Description string
	Type        string // compose and manifest
	Status      string
	Content     string
	Version     string
}
