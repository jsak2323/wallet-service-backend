package healthcheck

type HealthCheckRepository interface {
    GetAll() ([]HealthCheck, error)
    GetByRpcConfigId(rpcConfigId int) (HealthCheck, error)
    Create(healthCheck *HealthCheck) (error)
    Update(healthCheck *HealthCheck) (error)
}