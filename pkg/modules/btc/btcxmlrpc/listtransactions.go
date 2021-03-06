package btcxmlrpc

import (
	"context"
	"errors"

	// "encoding/json"

	// "github.com/mitchellh/mapstructure"
	// "github.com/elliotchance/phpserialize"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	// logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type ListTransactionsNodeXmlRpcRes struct {
	Content ListTransactionsNodeXmlRpcResStruct
}
type ListTransactionsNodeXmlRpcResStruct struct {
	Transactions string
}

// type TransactionRes struct {
//     Address             string      `json:"address"`
//     Category            string      `json:"category"`
//     Amount              float64     `json:"amount"`
//     Label               string      `json:"label"`
//     Vout                int         `json:"vout"`
//     Confirmations       int         `json:"confirmations"`
//     BlockHash           string      `json:"blockhash"`
//     BlockIndex          int         `json:"blockindex"`
//     BlockTime           int         `json:"blocktime"`
//     TxId                string      `json:"txid"`
//     WalletConflicts     []string    `json:"walletconflicts"`
//     Time                int         `json:"time"`
//     TimeReceived        int         `json:"timereceived"`
//     Bip125Replaceable   string      `json:"bip125-replaceable"`
//     Abandoned           bool        `json:"abandoned"`
// }

// todo: add limit compatibility
func (bs *BtcService) ListTransactions(ctx context.Context, rpcConfig rc.RpcConfig, limit int) (*model.ListTransactionsRpcRes, error) {
	res := model.ListTransactionsRpcRes{}

	rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
	xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

	nodeRpcRes := ListTransactionsNodeXmlRpcRes{}

	err := xmlrpc.XmlRpcCall("listtransactions", &rpcReq, &nodeRpcRes)

	if err == nil {
		// transactionsJson, err := serializedTransactionsToJson(nodeRpcRes.Content.Transactions)
		// if err != nil {
		//     logger.ErrorLog(" ---- btcxmlrpc ListTransactions serializedTransactionsToJson(nodeRpcRes.Content.Transactions) transactions: "+nodeRpcRes.Content.Transactions+", err: "+err.Error())
		//     return &res, err
		// }
		// res.Transactions = transactionsJson

		// res.Transactions = nodeRpcRes.Content.Transactions

		return &res, nil

	} else if err != nil {
		return &res, err

	} else {
		return &res, errors.New("Unexpected error occured in Node.")
	}
}

// func serializedTransactionsToJson(serializedTransactions string) (string, error) {
//     unserializedTransactions := []interface{}{}
//     phpserialize.Unmarshal([]byte(serializedTransactions), &unserializedTransactions)

//     transactions := []TransactionRes{}
//     for _, transaction := range unserializedTransactions {
//         transactionObj := TransactionRes{}
//         mapstructure.Decode(transaction, &transactionObj)
//         transactions = append(transactions, transactionObj)
//     }

//     transactionsJson, err := json.Marshal(transactions)
//     if err != nil { return "", err }

//     return string(transactionsJson), nil
// }
