package model

import (
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type GetBlockCountRpcRes struct {
    Blocks string
}

type ModuleService interface {
    GetSymbol() (string)
    GetHealthCheckRepo() (hc.HealthCheckRepository)
    GetBlockCount(rpcConfig rc.RpcConfig) (*GetBlockCountRpcRes, error)
    IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error)
}