package domain

import "context"

type Application struct {
	ID          uint
	Name        string
	Description string
	Type        string
	Content     string
	Version     string
}

type Repository interface {
	List(context.Context) ([]Application, error)
	Get(context.Context, uint) (Application, error)
	Delete(context.Context, uint) error
	Add(context.Context, *Application) error
	Update(context.Context, *Application) error
}
