package model

import (
	"fmt"

	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type GetBalanceRpcRes struct {
    Balance string
    Error   string
}

func (r *GetBalanceRpcRes) SetFromMapValues(mapValues map[string]interface{}, resFieldMap map[string]rrs.RpcResponse) (err error) {
	var ok bool
	var errRpcResp, balanceRpcResp rrs.RpcResponse

	if errRpcResp, ok = resFieldMap[rrs.FieldNameError]; !ok {
		return fmt.Errorf("Error rpc_response not configured")
	}
	
	// if error found, assign error to error field and return
	if ok := errRpcResp.ParseField(mapValues[rrs.FieldNameError], &r.Error); ok {
		return nil
	}

	if balanceRpcResp, ok = resFieldMap[rrs.FieldNameBalance]; !ok {
		return fmt.Errorf("Balance rpc_response not configured")
	}

	if ok = balanceRpcResp.ParseField(mapValues[rrs.FieldNameBalance], &r.Balance); !ok {
		return fmt.Errorf("misconfigured rpc_response id: %d", balanceRpcResp.Id)
	}

	return nil
}
