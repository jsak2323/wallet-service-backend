package model

import (
	"errors"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type GetNewAddressRpcRes struct {
    Address string
    Error   string
}

func (r *GetNewAddressRpcRes) SetFromMapValues(mapValues map[string]interface{}) (err error) {
	var ok bool

	if r.Address, ok = mapValues[rr.FieldNameAddress].(string); ok {
		return nil
	}

	// if not ok, look for error tag
	if r.Error, ok = mapValues[rr.FieldNameError].(string); !ok {
		return errors.New("mismatched rpc response data type")
	}

	return nil
}
