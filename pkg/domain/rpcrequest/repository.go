package rpcrequest

type Repository interface {
	Create(RpcRequest) error
	GetByRpcMethodId(rpcMethodId int) ([]RpcRequest, error)
	Update(RpcRequest) error
	Delete(id int) error
}
