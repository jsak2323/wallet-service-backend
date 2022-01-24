package rpcrequest

type Repository interface {
	Create(RpcRequest) error
	GetByRpcMethodId(rpcMethodId int) ([]RpcRequest, error)
	Update(UpdateRpcRequest) error
	Delete(id int) error
}
