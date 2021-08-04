package fireblocks

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

const createTxEndpoint = "createTransaction"
const VaultAccountType = "VAULT_ACCOUNT"
const InternalWalletType = "INTERNAL_WALLET"

func CreateTransaction(req CreateTransactionReq) (RES CreateTransactionRes, err error) {
	httpClient := &http.Client{
        Timeout: 120 * time.Second,
    }

	body, err := json.Marshal(req)
	if err != nil {
		return CreateTransactionRes{}, err
	}
	
	res, err := httpClient.Post(config.CONF.FireblocksHost+"/"+createTxEndpoint+"/", "application/json", bytes.NewBuffer(body))
    if err != nil {
        return CreateTransactionRes{}, err
    }
    defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&RES); err != nil {
		return CreateTransactionRes{}, err
	}

	return RES, nil
}