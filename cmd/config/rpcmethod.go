package config

import (
	"context"
	"errors"

	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
)

var RPCMETHOD = make(map[int]map[string]rm.RpcMethod) // map by rpcconfigid.rpcmethodtype

func LoadRpcMethodByRpcConfigId(ctx context.Context, rpcMethodRepo rm.Repository, rpcConfigId int) error {
	RPCMETHOD[rpcConfigId] = make(map[string]rm.RpcMethod)

	rpcMethods, err := rpcMethodRepo.GetByRpcConfigId(ctx, rpcConfigId)
	if err != nil {
		return err
	}

	for _, rpcMethod := range rpcMethods {
		RPCMETHOD[rpcConfigId][rpcMethod.Type] = rpcMethod
	}

	return nil
}

func GetRpcMethod(ctx context.Context, rpcMethodRepo rm.Repository, rpcConfgId int, rpcMethodType string) (rpcMethod rm.RpcMethod, err error) {
	if _, ok := RPCMETHOD[rpcConfgId]; ok {
		if _, ok := RPCMETHOD[rpcConfgId][rpcMethodType]; ok {
			return RPCMETHOD[rpcConfgId][rpcMethodType], nil
		}
	}

	if err = LoadRpcMethodByRpcConfigId(ctx, rpcMethodRepo, rpcConfgId); err != nil {
		return rm.RpcMethod{}, err
	}

	rpcMethod, ok := RPCMETHOD[rpcConfgId][rpcMethodType]
	if !ok {
		return rm.RpcMethod{}, errors.New("RPC method not found")
	}

	return rpcMethod, nil
}
