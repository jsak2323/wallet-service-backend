package model

import (
	"context"
	"errors"

	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

func InvalidRpcRequestConfig(name, method string) error {
	return errors.New("Invalid rpc request (name, method) -> (" + name + ")" + "(" + method + ")")
}

type AddressTypeRpcRes struct {
	AddressType string
	Error       string
}

type ModuleService interface {
	GetSymbol() string
	GetHealthCheckRepo() hc.Repository
	IsBlockCountHealthy(ctx context.Context, nodeBlockCount int, rpcConfigId int) (bool, int, error)

	GetBlockCount(ctx context.Context, rpcConfig rc.RpcConfig) (*GetBlockCountRpcRes, error)
	GetBalance(ctx context.Context, rpcConfig rc.RpcConfig) (*GetBalanceRpcRes, error)
	ListTransactions(ctx context.Context, rpcConfig rc.RpcConfig, limit int) (*ListTransactionsRpcRes, error)
	ListWithdraws(ctx context.Context, rpcConfig rc.RpcConfig, limit int) (*ListWithdrawsRpcRes, error)
	SendToAddress(ctx context.Context, rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*SendToAddressRpcRes, error)
	GetNewAddress(ctx context.Context, rpcConfig rc.RpcConfig, addressType string) (*GetNewAddressRpcRes, error)
	AddressType(ctx context.Context, rpcConfig rc.RpcConfig, address string) (*AddressTypeRpcRes, error)
}

type RpcRes interface {
	SetFromMapValues(mapValues map[string]interface{}, resFieldMap map[string]rrs.RpcResponse) (err error)
}
