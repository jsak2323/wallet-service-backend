package model

import (
	"fmt"
	"github.com/mitchellh/mapstructure"

	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type ListWithdrawsRpcRes struct {
	Withdraws []Withdraw
	Error     string
}

type Withdraw struct {
	AddressTo     string `json:"to"`
	Tx            string `json:"hash"`
	Amount        string `json:"amount"`
	Confirmations int    `json:"confirmations"`
	BlockchainFee string `json:"blockchain_fee"`
	Memo          string `json:"memo"`
	SuccessTime   string `json:"success_time"`
}

func (r *ListWithdrawsRpcRes) SetFromMapValues(mapValues map[string]interface{}, resFieldMap map[string]rrs.RpcResponse) (err error) {
	var ok bool
	var errRpcResp, transactionRpcresp rrs.RpcResponse

	if errRpcResp, ok = resFieldMap[rrs.FieldNameError]; !ok {
		return fmt.Errorf("Error rpc_response not configured")
	}

	// if error found, assign error to error field and return
	if ok = errRpcResp.ParseField(mapValues[rrs.FieldNameError], &r.Error); ok {
		return nil
	}

	if transactionRpcresp, ok = resFieldMap[rrs.FieldNameTransactions]; !ok {
		return fmt.Errorf("Error rpc_response not configured")
	}

	tempResult, err := transactionRpcresp.ParseArrayOfJson(mapValues[rrs.FieldNameTransactions])
	if err != nil {
		return err
	}

	for _, temp := range tempResult {
		transaction := Withdraw{}

		if err = mapstructure.Decode(temp, &transaction); err != nil {
			return err
		}

		r.Withdraws = append(r.Withdraws, transaction)
	}

	return nil
}
