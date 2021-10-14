package cold

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *ColdWalletService) UpdateBalanceHandler(w http.ResponseWriter, req *http.Request) {
	var updateReq UpdateBalanceReq
	var RES StandardRes
	var err error

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
		logger.ErrorLog(" - UpdateBalance json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = s.UpdateBalance(updateReq.Id, updateReq.Balance); err != nil {
		logger.ErrorLog(" - UpdateBalanceHandler s.UpdateBalance err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func (s *ColdWalletService) UpdateBalance(id int, balance string) error {
	balanceRaw, err := util.CoinToRaw(balance, 8)
	if err != nil {
		return err
	}

	if err = s.cbRepo.UpdateBalance(id, balanceRaw); err != nil {
		return err
	}

	return nil
}
