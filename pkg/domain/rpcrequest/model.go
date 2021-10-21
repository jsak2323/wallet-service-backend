package rpcrequest

type RpcRequest struct {
	Id          int
	ArgName     string
	ArgOrder    int
	Source      string
	Value       string
	RpcMethodId int
}
