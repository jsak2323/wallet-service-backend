package healthcheck

type HealthCheck struct {
    Id                      int
    RpcConfigId             int
    BlockCount              int
    ConfirmBlockCount       int
    LastUpdated             string
}