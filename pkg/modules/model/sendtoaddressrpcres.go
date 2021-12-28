package model

import (
	"errors"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type SendToAddressRpcRes struct {
    TxHash  string
    Error   string
}

func (r *SendToAddressRpcRes) SetFromMapValues(mapValues map[string]interface{}, resFieldMap map[string]rrs.RpcResponse) (err error) {
	var ok bool

	if r.TxHash, ok = mapValues[rr.FieldNameTxHash].(string); ok {
		return nil
	}

	// if not ok, look for error tag
	if r.Error, ok = mapValues[rr.FieldNameError].(string); !ok {
		return errors.New("mismatched rpc response data type")
	}

	return nil
}
