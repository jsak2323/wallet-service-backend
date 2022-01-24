package rpcconfig

import (
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type GetRes struct {
	RpcConfig domain.RpcConfig `json:"rpc_config"`
	Error     *errs.Error      `json:"error"`
}

type ListRes struct {
	RpcConfigs []domain.RpcConfig `json:"rpc_configs"`
	Error      *errs.Error        `json:"error"`
}

type RpcConfigRpcMethodReq struct {
	RpcConfigId int `json:"rpc_config_id" validate:"required"`
	RpcMethodId int `json:"rpc_method_id" validate:"required"`
}

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}
