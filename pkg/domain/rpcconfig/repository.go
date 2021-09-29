package rpcconfig

import (
	"errors"
)

type RpcConfigRepository interface {
    GetByCurrencyId(currency_id int) ([]RpcConfig, error)
    GetByCurrencySymbol(symbol string) ([]RpcConfig, error)
}

const ReceiverRpcType = "receiver"
const MasterRpcType = "master"

var errReceiverNotFound = errors.New("no receiver rpc found")


func GetReceiverFromList(rcs []RpcConfig) (rc RpcConfig, err error) {
	for _, rc := range rcs {
		if rc.Type == ReceiverRpcType || rc.Type == MasterRpcType {
			return rc, nil
		}
	}

	return rc, errReceiverNotFound
}
