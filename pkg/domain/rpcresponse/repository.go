package rpcresponse

type Repository interface {
	Create(CreateRpcResponse) error
	GetByRpcMethodId(rpcMethodId int) ([]RpcResponse, error)
	Update(RpcResponse) error
	Delete(id int) error
}
