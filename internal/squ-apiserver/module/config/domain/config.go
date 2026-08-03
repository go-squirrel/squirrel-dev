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
	Add(context.Context, *Config) error
	Update(context.Context, *Config) error
}
