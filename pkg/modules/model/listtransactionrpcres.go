package model

import (
	"errors"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type ListTransactionsRpcRes struct {
    Transactions string
    Error        string
}

func (r *ListTransactionsRpcRes) SetFromMapValues(mapValues map[string]interface{}) (err error) {
	var ok bool

	if r.Transactions, ok = mapValues[rr.FieldNameBlockCount].(string); ok {
		return nil
	}

	// if not ok, look for error tag
	if r.Error, ok = mapValues[rr.FieldNameError].(string); !ok {
		return errors.New("mismatched rpc response data type")
	}

	return nil
}
