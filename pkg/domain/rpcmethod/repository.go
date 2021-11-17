package rpcmethod

type Repository interface {
	Create(RpcMethod) (int, error)
	GetAll(page, limit int) ([]RpcMethod, error)
	GetByRpcConfigId(rpcConfigId int) ([]RpcMethod, error)
	Update(RpcMethod) error
	Delete(id int) error
}
