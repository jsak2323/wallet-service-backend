package rpcmethod

import "context"

type Repository interface {
	Create(context.Context, RpcMethod) (int, error)
	GetAll(ctx context.Context, page, limit int) ([]RpcMethod, error)
	GetByRpcConfigId(ctx context.Context, rpcConfigId int) ([]RpcMethod, error)
	Update(context.Context, UpdateRpcMethod) error
	Delete(id int) error
}
