package api

import (
	"squirrel-dev/internal/squ-apiserver/module/config/api/res"
	"squirrel-dev/internal/squ-apiserver/module/config/domain"
)

func fromDomain(value domain.Config) res.Config {
	return res.Config{
		ID:    value.ID,
		Key:   value.Key,
		Value: value.Value,
	}
}
