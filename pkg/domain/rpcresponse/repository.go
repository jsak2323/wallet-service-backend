package rpcresponse

import "context"

type Repository interface {
	Create(context.Context, CreateRpcResponse) error
	GetByRpcMethodId(ctx context.Context, rpcMethodId int) ([]RpcResponse, error)
	Update(context.Context, RpcResponse) error
	Delete(ctx context.Context, id int) error
}
