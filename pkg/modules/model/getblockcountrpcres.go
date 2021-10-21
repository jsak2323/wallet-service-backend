package model

import (
	"errors"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type GetBlockCountRpcRes struct {
    Blocks  string
    Error   string
}

func (r *GetBlockCountRpcRes) SetFromMapValues(mapValues map[string]interface{}) (err error) {
	var ok bool

	if r.Blocks, ok = mapValues[rr.FieldNameBlockCount].(string); ok {
		return nil
	}

	// if not ok, look for error tag
	if r.Error, ok = mapValues[rr.FieldNameError].(string); !ok {
		return errors.New("mismatched rpc response data type")
	}

	return nil
}
