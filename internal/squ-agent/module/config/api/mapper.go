package api

import (
	"squirrel-dev/internal/squ-agent/module/config/api/res"
	"squirrel-dev/internal/squ-agent/module/config/domain"
)

func fromDomain(value domain.Config) res.Config {
	return res.Config{
		ID:    value.ID,
		Key:   value.Key,
		Value: value.Value,
	}
}
