package infra

import (
	"context"

	configDomain "squirrel-dev/internal/squ-agent/module/config/domain"
)

type ConfigStore struct {
	repository configDomain.Repository
}

func NewConfigStore(repository configDomain.Repository) *ConfigStore {
	return &ConfigStore{repository: repository}
}

func (s *ConfigStore) Save(ctx context.Context, key, value string) error {
	return s.repository.CreateOrUpdate(ctx, &configDomain.Config{Key: key, Value: value})
}
