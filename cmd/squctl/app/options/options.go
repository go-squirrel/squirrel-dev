package options

import (
	"squirrel-dev/internal/squctl/config"
	"squirrel-dev/internal/squctl/server"
)

type AppOptions struct {
	ConfFile string
	Config   *config.Config
}

func NewAppOptions() *AppOptions {
	o := &AppOptions{}
	return o
}

func (o *AppOptions) NewServer() (*server.Server, error) {
	s := server.NewServer()
	o.loadConfig(o.ConfFile)
	s.Config = o.Config

	return s, nil
}

func (o *AppOptions) loadConfig(configFile string) {
	o.Config = config.New(configFile)
}
