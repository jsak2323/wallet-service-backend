package rpcmethod

type RpcMethod struct {
	Id          int
	Name        string `validate:"required"`
	Type        string `validate:"required"`
	NumOfArgs   int
	Network     string
	RpcConfigId int `validate:"required"`
}

type UpdateRpcMethod struct {
	Id          int    `validate:"required"`
	Name        string `validate:"required"`
	Type        string `validate:"required"`
	NumOfArgs   int
	Network     string
	RpcConfigId int
}
