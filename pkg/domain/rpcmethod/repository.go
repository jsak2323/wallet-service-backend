package rpcmethod

type Repository interface {
	GetByRpcConfigId(rpcConfigId int) ([]RpcMethod, error)
}
