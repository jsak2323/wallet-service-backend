package currencyrpc

import "context"

type Repository interface {
	Create(ctx context.Context, currencyConfigId, RpcConfigId int) error
	GetByCurrencyConfig(currencyConfigId int) ([]CurrencyRpc, error)
	GetByRpcConfig(RpcConfigId int) ([]CurrencyRpc, error)
	DeleteByCurrencyConfigId(currencyConfigId int) error
	DeleteByRpcConfigId(RpcConfigId int) error
	Delete(ctx context.Context, currencyConfigId, RpcConfigId int) error
}
