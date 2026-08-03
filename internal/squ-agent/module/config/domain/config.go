package domain

import "context"

type Config struct {
	ID    uint
	Key   string
	Value string
}

type Repository interface {
	List(context.Context) ([]Config, error)
	Get(context.Context, uint) (Config, error)
	Delete(context.Context, uint) error
	CreateOrUpdate(context.Context, *Config) error
	GetByKey(context.Context, string) (string, error)
	Transaction(context.Context, func(Repository) error) error
}
