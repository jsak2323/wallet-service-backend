package modules

import (
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
)

type GetBlockCountRpcRes struct {
    Blocks string
}

type ModuleService interface {
    GetBlockCount(rpcConfig rc.RpcConfig) (*GetBlockCountRpcRes, error)
}

func NewModuleServices() *map[string]ModuleService {
    ModuleServices := make(map[string]ModuleService)

    ModuleServices["BTC"] = &BtcService{}
    ModuleServices["ETH"] = &EthService{}

    return &ModuleServices
}