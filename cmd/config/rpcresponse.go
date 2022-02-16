package config

import (
	"context"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

var RPCRRESPONSE = make(map[int]map[string]rr.RpcResponse) // map by rpcmethodid.rpcresponsefieldname

func LoadRpcResponseByRpcMethodId(ctx context.Context, rpcResponseRepo rr.Repository, rpcMethodId int) error {
	RPCRRESPONSE[rpcMethodId] = make(map[string]rr.RpcResponse)

	rpcResponses, err := rpcResponseRepo.GetByRpcMethodId(ctx, rpcMethodId)
	if err != nil {
		return err
	}

	for _, rpcResponse := range rpcResponses {
		RPCRRESPONSE[rpcMethodId][rpcResponse.TargetFieldName] = rpcResponse
	}

	return nil
}

func GetRpcResponseMap(ctx context.Context, rpcResponseRepo rr.Repository, rpcMethodId int) (rpcResponse map[string]rr.RpcResponse, err error) {
	if _, ok := RPCRRESPONSE[rpcMethodId]; ok {
		return RPCRRESPONSE[rpcMethodId], nil
	}

	if err = LoadRpcResponseByRpcMethodId(ctx, rpcResponseRepo, rpcMethodId); err != nil {
		return map[string]rr.RpcResponse{}, err
	}

	return RPCRRESPONSE[rpcMethodId], nil
}
