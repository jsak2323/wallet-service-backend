package rpcconfig

import (
	"errors"
)

type RpcConfigRepository interface {
	GetByCurrencyId(currency_id int) ([]RpcConfig, error)
	GetByCurrencySymbol(symbol string) ([]RpcConfig, error)
}

const MasterRpcType = "master"
const SenderRpcType = "sender"
const ReceiverRpcType = "receiver"

var errSenderNotFound = errors.New("no sender rpc found")
var errReceiverNotFound = errors.New("no receiver rpc found")

func GetSenderFromList(rcs []RpcConfig) (rc RpcConfig, err error) {
	for _, rc := range rcs {
		if rc.Type == "sender" || rc.Type == "master" {
			return rc, nil
		}
	}

	return rc, errSenderNotFound
}

func GetReceiverFromList(rcs []RpcConfig) (rc RpcConfig, err error) {
	for _, rc := range rcs {
		if rc.Type == ReceiverRpcType || rc.Type == MasterRpcType {
			return rc, nil
		}
	}

	return rc, errReceiverNotFound
}
