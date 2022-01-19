package cold

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *ColdWalletService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var updateReq domain.ColdBalance
	var RES StandardRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
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

	if err = validateCreateReq(updateReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	updateReq.Balance = util.CoinToRaw(updateReq.Balance, 8)

	if err = s.cbRepo.Update(updateReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateColdBalance)
		return
	}
}

func validateUpdateReq(updateReq domain.ColdBalance) error {
	if updateReq.Id == 0 {
		return errs.AddTrace(errors.New("Id"))
	}
	if updateReq.CurrencyId == 0 {
		return errs.AddTrace(errors.New("Currency Id"))
	}
	if updateReq.Name == "" {
		return errs.AddTrace(errors.New("Name"))
	}
	if updateReq.Type == "" {
		return errs.AddTrace(errors.New("Type"))
	}
	if updateReq.Address == "" {
		return errs.AddTrace(errors.New("Address"))
	}

	if isFireblocksCold(updateReq.Type) && updateReq.FireblocksName == "" {
		return errs.AddTrace(errors.New("Fireblocks Name"))
	}

	return nil
}
