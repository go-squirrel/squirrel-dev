package domain

import "context"

type App struct {
	ID          uint
	Name        string
	Description string
	Type        string
	Category    string
	Icon        *string
	Version     string
	Content     string
	Tags        string
	Author      string
	RepoURL     *string
	HomepageURL *string
	IsOfficial  bool
	Downloads   int
	Status      string
}

type Repository interface {
	List(context.Context) ([]App, error)
	Get(context.Context, uint) (App, error)
	Delete(context.Context, uint) error
	Add(context.Context, *App) error
	Update(context.Context, *App) error
}
