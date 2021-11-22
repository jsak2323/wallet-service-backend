package rpcconfigrpcmethod

type Repository interface {
	Create(rpcConfigId, rpcMethodId int) error
	GetByRpcConfig(rpcConfigId int) ([]RpcConfigRpcMethod, error)
	GetByRpcMethod(rpcMethodId int) ([]RpcConfigRpcMethod, error)
	DeleteByRpcConfig(rpcConfigId int) error
	DeleteByRpcMethod(rpcMethodId int) error
	Delete(rpcConfigId, rpcMethodId int) error
}
