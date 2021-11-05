package rpcrequest

func JsonFieldTypeRpcRequests(rpcRequests []RpcRequest) (result []RpcRequest) {
	for _, field := range rpcRequests {
		if field.Type == TypeJsonField {
			result = append(result, field)
		}
	}

	return result
}
