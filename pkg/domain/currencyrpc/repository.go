package currencyrpc

type Repository interface {
	Create(currencyConfigId, RpcConfigId int) error
	GetByCurrencyConfig(currencyConfigId int) ([]CurrencyRpc, error)
	GetByRpcConfig(RpcConfigId int) ([]CurrencyRpc, error)
	DeleteByCurrencyConfigId(currencyConfigId int) error
	DeleteByRpcConfigId(RpcConfigId int) error
	Delete(currencyConfigId, RpcConfigId int) error
}
