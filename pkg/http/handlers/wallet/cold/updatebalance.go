package cold

import (
	"encoding/json"
	"net/http"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (s *ColdWalletService) UpdateBalanceHandler(w http.ResponseWriter, req *http.Request) {
	var updateReq UpdateBalanceReq
	var RES StandardRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Cold balance successfully updated"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&updateReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = s.UpdateBalance(updateReq.Id, updateReq.Balance); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateColdBalance)
		return
	}
}

func (s *ColdWalletService) UpdateBalance(id int, balance string) (err error) {
	balanceRaw := util.CoinToRaw(balance, 8)

	if err = s.cbRepo.UpdateBalance(id, balanceRaw); err != nil {
		return errs.AddTrace(err)
	}

	return nil
}
