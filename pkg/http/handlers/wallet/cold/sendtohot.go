package cold

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
)

func (s *ColdWalletService) SendToHotHandler(w http.ResponseWriter, req *http.Request) {
	var sendToHotReq SendToHotReq
	var RES StandardRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Cold balance successfully sent"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&sendToHotReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	vaultAccountId, err := FireblocksVaultAccountId(sendToHotReq.FireblocksType)
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedFireblocksVaultAccountId)
		return
	}

	res, err := fireblocks.CreateTransaction(fireblocks.CreateTransactionReq{
		AssetId: sendToHotReq.FireblocksName,
		Amount:  sendToHotReq.Amount,
		Source: fireblocks.TransactionAccount{
			Type: fireblocks.VaultAccountType,
			Id:   vaultAccountId,
		},
		Destination: fireblocks.TransactionAccount{
			Type: fireblocks.InternalWalletType,
			Id:   config.CONF.FireblocksHotVaultId,
		},
	})
	if err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateTransaction)
		return
	}

	if res.Error != "" {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateTransaction)
		return
	}
}
