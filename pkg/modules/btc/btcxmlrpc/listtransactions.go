package btcxmlrpc

import(
    "errors"
    // "encoding/json"

    // "github.com/mitchellh/mapstructure"
    // "github.com/elliotchance/phpserialize"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    // logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type ListTransactionsNodeXmlRpcRes struct {
    Response ListTransactionsNodeXmlRpcResStruct
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
func (bs *BtcService) ListTransactions(rpcConfig rc.RpcConfig, limit int) (*model.ListTransactionsRpcRes, error) {
    res := model.ListTransactionsRpcRes{ Transactions: "" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    nodeRpcRes := ListTransactionsNodeXmlRpcRes{}

    err := xmlrpc.XmlRpcCall("listtransactions", &rpcReq, &nodeRpcRes)

    if err == nil {
        // transactionsJson, err := serializedTransactionsToJson(nodeRpcRes.Response.Transactions)
        // if err != nil {
        //     logger.ErrorLog(" ---- btcxmlrpc ListTransactions serializedTransactionsToJson(nodeRpcRes.Response.Transactions) transactions: "+nodeRpcRes.Response.Transactions+", err: "+err.Error())
        //     return &res, err
        // }
        // res.Transactions = transactionsJson

        res.Transactions = nodeRpcRes.Response.Transactions

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