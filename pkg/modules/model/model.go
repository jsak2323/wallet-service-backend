package model

import (
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
)

type GetBlockCountRpcRes struct {
    Blocks string
}

type ModuleService interface {
    GetBlockCount(rpcConfig rc.RpcConfig) (*GetBlockCountRpcRes, error)
}