package model

import (
    "errors"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
)

func InvalidRpcRequestConfig(name, method string) error {
    return errors.New("Invalid rpc request (name, method) -> (" + name + ")" + "(" + method + ")")
}

type AddressTypeRpcRes struct {
    AddressType string
    Error       string
}

type ModuleService interface {
    GetSymbol() (string)
    GetHealthCheckRepo() (hc.Repository)
    IsBlockCountHealthy(nodeBlockCount int, rpcConfigId int) (bool, int, error)

    GetBlockCount(rpcConfig rc.RpcConfig) (*GetBlockCountRpcRes, error)
    GetBalance(rpcConfig rc.RpcConfig) (*GetBalanceRpcRes, error)
    ListTransactions(rpcConfig rc.RpcConfig, limit int) (*ListTransactionsRpcRes, error)
    ListWithdraws(rpcConfig rc.RpcConfig, limit int) (*ListWithdrawsRpcRes, error)
    SendToAddress(rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*SendToAddressRpcRes, error)
    GetNewAddress(rpcConfig rc.RpcConfig, addressType string) (*GetNewAddressRpcRes, error)
    AddressType(rpcConfig rc.RpcConfig, address string) (*AddressTypeRpcRes, error)
}

type RpcRes interface {
    SetFromMapValues(mapValues map[string]interface{}, resFieldMap map[string]rrs.RpcResponse) (err error)
}

