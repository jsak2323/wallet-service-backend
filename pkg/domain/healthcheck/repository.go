package healthcheck

type HealthCheckRepository interface {
    GetAll() ([]BlockCount, error)
}