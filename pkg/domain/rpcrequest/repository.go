package rpcrequest

import "context"

type Repository interface {
	Create(context.Context, RpcRequest) error
	GetByRpcMethodId(rpcMethodId int) ([]RpcRequest, error)
	Update(UpdateRpcRequest) error
	Delete(id int) error
}
