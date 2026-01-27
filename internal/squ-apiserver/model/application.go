package model

type Application struct {
	BaseModel
	Name        string
	Description string
	Type        string // compose and manifest
	Content     string
	Version     string
}
