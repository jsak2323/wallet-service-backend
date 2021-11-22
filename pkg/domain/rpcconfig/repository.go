package rpcconfig

type Repository interface {
	Create(RpcConfig) error
	GetAll(page, limit int) ([]RpcConfig, error)
	GetById(id int) (RpcConfig, error)
	GetByCurrencyId(currency_id int) ([]RpcConfig, error)
	GetByCurrencySymbol(symbol string) ([]RpcConfig, error)
	Update(RpcConfig) error
	ToggleActive(Id int, active bool) error
}
