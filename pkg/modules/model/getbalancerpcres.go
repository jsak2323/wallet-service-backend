package model

import (
	"errors"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type GetBalanceRpcRes struct {
    Balance string
    Error   string
}

func (r *GetBalanceRpcRes) SetFromMapValues(mapValues map[string]interface{}) (err error) {
	var ok bool

	if r.Balance, ok = mapValues[rr.FieldNameBalance].(string); ok {
		return nil
	}

	// if not ok, look for error tag
	if r.Error, ok = mapValues[rr.FieldNameError].(string); !ok {
		return errors.New("mismatched rpc response data type")
	}

	return nil
}
