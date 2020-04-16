package model

import (
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type GetBlockCountRpcRes struct {
    Blocks string
}

type GetBalanceRpcRes struct {
    Balance string
}

type ListTransactionsRpcRes struct {
    Transactions string
}

type SendToAddressRpcRes struct {
    TxHash string
}

type ModuleService interface {
    GetSymbol() (string)
    GetHealthCheckRepo() (hc.HealthCheckRepository)
    IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error)

    GetBlockCount(rpcConfig rc.RpcConfig) (*GetBlockCountRpcRes, error)
    GetBalance(rpcConfig rc.RpcConfig) (*GetBalanceRpcRes, error)
    ListTransactions(rpcConfig rc.RpcConfig, limit int) (*ListTransactionsRpcRes, error)
    SendToAddress(rpcConfig rc.RpcConfig, address string, amountInDecimal string) (*SendToAddressRpcRes, error)
}