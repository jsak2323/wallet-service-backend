package healthcheck

type HealthCheckRepository interface {
    GetAll() ([]HealthCheck, error)
}