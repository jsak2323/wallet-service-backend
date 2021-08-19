package fireblocks

import (
	"errors"
	"encoding/json"
	
	"gopkg.in/resty.v0"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

const createTxEndpoint = "createTransaction"
const VaultAccountType = "VAULT_ACCOUNT"
const InternalWalletType = "INTERNAL_WALLET"

func CreateTransaction(req CreateTransactionReq) (RES CreateTransactionRes, err error) {
	body, err := json.Marshal(req)
	if err != nil {
		return CreateTransactionRes{}, err
	}
	
	res, err := resty.R().
		SetHeader("Authorization", "Basic " + auth()).
		SetBody(body).
		Post(config.CONF.FireblocksHost+"/"+createTxEndpoint+"/")

    if err != nil {
        return CreateTransactionRes{}, err
    }

	if err = json.Unmarshal(res.Body, &RES); err != nil {
		return CreateTransactionRes{}, err
	}

	if RES.Error != "" {
		return CreateTransactionRes{}, errors.New(RES.Error)
	}

	return RES, nil
}