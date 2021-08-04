package cold

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *ColdWalletService) UpdateBalanceHandler(w http.ResponseWriter, req *http.Request) {
	var updateReq UpdateReq
	var RES       StandardRes
	var err       error

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Cold balance successfully updated"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		logger.ErrorLog(" - UpdateBalanceHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	balanceRaw, err := util.CoinToRaw(updateReq.Balance, 8)
	if err != nil {
		logger.ErrorLog(" - UpdateBalanceHandler CoinToRaw err: " + err.Error())
		return
	}

	if err = s.cbRepo.UpdateBalance(updateReq.Id, balanceRaw); err != nil {
		logger.ErrorLog(" - UpdateBalanceHandler cbRepo.UpdateBalance err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}