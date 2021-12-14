package model

import (
	"fmt"
	"github.com/mitchellh/mapstructure"

	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
)

type ListTransactionsRpcRes struct {
	Transactions []Transaction
	Error        string
}

type Transaction struct {
	AddressTo   string `json:"to"`
	Tx          string `json:"hash"`
	Amount      string `json:"amount"`
	Memo        string `json:"memo"`
	SuccessTime string `json:"success_time"`
}

func (r *ListTransactionsRpcRes) SetFromMapValues(mapValues map[string]interface{}, resFieldMap map[string]rrs.RpcResponse) (err error) {
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
		transaction := Transaction{}

		if err = mapstructure.Decode(temp, &transaction); err != nil {
			return err
		}

		r.Transactions = append(r.Transactions, transaction)
	}

	return nil
}
