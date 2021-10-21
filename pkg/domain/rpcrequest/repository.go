package rpcrequest

type Repository interface {
	GetByRpcMethodId(rpcMethodId int) ([]RpcRequest, error)
}
