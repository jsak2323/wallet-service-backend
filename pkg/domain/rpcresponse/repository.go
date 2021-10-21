package rpcresponse

type Repository interface {
	GetByRpcMethodId(rpcMethodId int) ([]RpcResponse, error)
}
