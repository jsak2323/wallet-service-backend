package rpcconfig

import (
	"errors"
)

const MasterRpcType string = "master"
const SenderRpcType string = "sender"
const ReceiverRpcType string = "receiver"

var errSenderNotFound error = errors.New("no sender rpc found")
var errReceiverNotFound error = errors.New("no receiver rpc found")
