package cold

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *ColdWalletService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var updateReq domain.ColdBalance
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
		logger.ErrorLog(" - UpdateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateCreateReq(updateReq); err != nil {
		logger.ErrorLog(" - CreateHandler validateCreateReq err: " + err.Error())
		RES.Error = "Invalid param: " + err.Error()
		return
	}

	if updateReq.Balance, err = util.CoinToRaw(updateReq.Balance, 8); err != nil {
		logger.ErrorLog(" - UpdateHandler util.CoinToRaw err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = s.cbRepo.Update(updateReq); err != nil {
		logger.ErrorLog(" - UpdateHandler s.cbRepo.Update err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateUpdateReq(updateReq domain.ColdBalance) error {
	if updateReq.Id == 0 {
		return errors.New("Id")
	}
	if updateReq.CurrencyId == 0 {
		return errors.New("Currency Id")
	}
	if updateReq.Name == "" {
		return errors.New("Name")
	}
	if updateReq.Type == "" {
		return errors.New("Type")
	}
	if updateReq.Address == "" {
		return errors.New("Address")
	}

	if isFireblocksCold(updateReq.Type) && updateReq.FireblocksName == "" {
		return errors.New("Fireblocks Name")
	}

	return nil
}
