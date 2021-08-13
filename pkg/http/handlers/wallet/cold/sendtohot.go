package cold

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *ColdWalletService) SendToHotHandler(w http.ResponseWriter, req *http.Request) {
	var sendToHotReq SendToHotReq
	var RES       	 StandardRes
	var err       	 error

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Cold balance successfully sent"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&sendToHotReq); err != nil {
		logger.ErrorLog(" - SendToHotHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	res, err := fireblocks.CreateTransaction(fireblocks.CreateTransactionReq{
		AssetId: sendToHotReq.FireblocksName,
		Amount: sendToHotReq.Amount,
		Source: fireblocks.TransactionAccount{
			Type: fireblocks.VaultAccountType,
			Id: FireblocksVaultAccountId(sendToHotReq.FireblocksType),
		},
		Destination: fireblocks.TransactionAccount{
			Type: fireblocks.InternalWalletType,
			Id: config.CONF.FireblocksHotVaultId,
		},
	})
	if err != nil {
		logger.ErrorLog(" - SendToHotHandler fireblocks.CreateTransaction err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if res.Error != "" {
		logger.ErrorLog(" - SendToHotHandler fireblocks.CreateTransaction err: " + res.Error)
		RES.Error = errInternalServer
		return
	}
}