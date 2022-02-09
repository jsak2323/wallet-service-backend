package rpcconfig

import "context"

type Repository interface {
	Create(context.Context, RpcConfig) error
	GetAll(ctx context.Context, page, limit int) ([]RpcConfig, error)
	GetById(ctx context.Context, id int) (RpcConfig, error)
	GetByCurrencyId(ctx context.Context, currency_id int) ([]RpcConfig, error)
	GetByCurrencySymbol(symbol string) ([]RpcConfig, error)
	Update(context.Context, UpdateRpcConfig) error
	ToggleActive(ctx context.Context, Id int, active bool) error
}
