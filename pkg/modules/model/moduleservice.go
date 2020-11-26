package model

import (
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type GetBlockCountRpcRes struct {
    Blocks  string
    Error   string
}

type GetBalanceRpcRes struct {
    Balance string
    Error   string
}

type ListTransactionsRpcRes struct {
    Transactions string
    Error        string
}

type SendToAddressRpcRes struct {
    TxHash  string
    Error   string
}

type GetNewAddressRpcRes struct {
    Address string
    Error   string
}

type AddressTypeRpcRes struct {
    AddressType string
    Error       string
}

type ModuleService interface {
    GetSymbol() (string)
    GetHealthCheckRepo() (hc.HealthCheckRepository)
    IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error)

    GetBlockCount(rpcConfig rc.RpcConfig) (*GetBlockCountRpcRes, error)
    GetBalance(rpcConfig rc.RpcConfig) (*GetBalanceRpcRes, error)
    ListTransactions(rpcConfig rc.RpcConfig, limit int) (*ListTransactionsRpcRes, error)
    SendToAddress(rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*SendToAddressRpcRes, error)
    GetNewAddress(rpcConfig rc.RpcConfig, addressType string) (*GetNewAddressRpcRes, error)
    AddressType(rpcConfig rc.RpcConfig, address string) (*AddressTypeRpcRes, error)
}


