package rpcconfig

type RpcConfigRepository interface {
	Create(RpcConfig) error
	GetById(id int) (RpcConfig, error)
	GetByCurrencyId(currency_id int) ([]RpcConfig, error)
	GetByCurrencySymbol(symbol string) ([]RpcConfig, error)
	Update(RpcConfig) error
	ToggleActive(Id int, active bool) error
}
