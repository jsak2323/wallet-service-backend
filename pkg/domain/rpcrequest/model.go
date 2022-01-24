package rpcrequest

type RpcRequest struct {
	Id             int
	ArgName        string `validate:"required"`
	Type           string `validate:"required"`
	ArgOrder       int
	Source         string `validate:"required"`
	RuntimeVarName string
	Value          string
	RpcMethodId    int `validate:"required"`
}

type UpdateRpcRequest struct {
	Id             int    `validate:"required"`
	ArgName        string `validate:"required"`
	Type           string `validate:"required"`
	ArgOrder       int
	Source         string `validate:"required"`
	RuntimeVarName string
	Value          string
	RpcMethodId    int `validate:"required"`
}
