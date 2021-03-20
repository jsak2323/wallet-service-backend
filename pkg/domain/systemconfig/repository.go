package systemconfig

type SystemConfigRepository interface {
    GetAll() ([]SystemConfig, error)
    GetByName(configName string) (*SystemConfig, error)
}


