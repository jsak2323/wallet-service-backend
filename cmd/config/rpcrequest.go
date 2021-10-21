package config

import (
	rr"github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
)

var RPCREQUEST = make(map[int][]rr.RpcRequest) // map by rpcmethodid.rpcrequest_argorder

func LoadRpcRequestByRpcMethodId(rpcRequestRepo rr.Repository, rpcMethodId int) (err error) {
	RPCREQUEST[rpcMethodId] = []rr.RpcRequest{}

	RPCREQUEST[rpcMethodId], err = rpcRequestRepo.GetByRpcMethodId(rpcMethodId)
	if err != nil {
		return err
	}

	return nil
}

func GetRpcRequestMap(rpcRequestRepo rr.Repository, rpcMethodId int) (rpcRequest []rr.RpcRequest, err error) {
	if _, ok := RPCREQUEST[rpcMethodId]; ok {
		return RPCREQUEST[rpcMethodId], nil
	}

	if err = LoadRpcRequestByRpcMethodId(rpcRequestRepo, rpcMethodId); err != nil {
		return []rr.RpcRequest{}, err
	}

	return RPCREQUEST[rpcMethodId], nil
}
