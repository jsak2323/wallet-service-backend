package rpcrequest

type RpcRequest struct {
	Id             int
	ArgName        string
	Type           string
	ArgOrder       int
	Source         string
	RuntimeVarName string
	Value          string
	RpcMethodId    int
}
