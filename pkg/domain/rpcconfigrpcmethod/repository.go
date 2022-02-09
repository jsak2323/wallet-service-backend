package rpcconfigrpcmethod

import "context"

type Repository interface {
	Create(ctx context.Context, rpcConfigId, rpcMethodId int) error
	GetByRpcConfig(rpcConfigId int) ([]RpcConfigRpcMethod, error)
	GetByRpcMethod(rpcMethodId int) ([]RpcConfigRpcMethod, error)
	DeleteByRpcConfig(rpcConfigId int) error
	DeleteByRpcMethod(ctx context.Context, rpcMethodId int) error
	Delete(ctx context.Context, rpcConfigId, rpcMethodId int) error
}
