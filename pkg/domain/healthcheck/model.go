package healthcheck

import (
    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
)

type HealthCheck struct {
    Id           int
    RpcConfigId  int
    BlockCount   int
    BlockDiff    int
    IsHealthy    bool
    LastUpdated  string
}

type HealthCheckWithRpcConfig struct {
    Id           int
    BlockCount   int
    BlockDiff    int
    IsHealthy    bool
    LastUpdated  string
    RpcConfig    rc.RpcConfig
}


