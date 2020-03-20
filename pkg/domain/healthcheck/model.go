package healthcheck

type HealthCheck struct {
    Id                      int
    RpcConfigId             int
    BlockCount              int
    BlockDiff               int
    IsHealthy               bool
    LastUpdated             string
}