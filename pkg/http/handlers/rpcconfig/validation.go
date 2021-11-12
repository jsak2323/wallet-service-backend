package rpcconfig

func (req RpcConfigRpcMethodReq) valid() bool {
	if req.RpcConfigId == 0 {
		return false
	}
	
	if req.RpcMethodId == 0 {
		return false
	}

	return true
}