package systemconfig

import "context"

type Repository interface {
	GetAll() ([]SystemConfig, error)
	GetByName(ctx context.Context, configName string) (*SystemConfig, error)
	Update(sysConf SystemConfig) error
}
