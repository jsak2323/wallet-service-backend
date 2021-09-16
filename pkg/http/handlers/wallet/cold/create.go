package cold

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *ColdWalletService) CreateHandler(w http.ResponseWriter, req *http.Request) {
	var createReq domain.ColdBalance
	var RES StandardRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Cold wallet successfully created"
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		logger.ErrorLog(" - CreateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateCreateReq(createReq); err != nil {
		logger.ErrorLog(" - CreateHandler validateCreateReq err: " + err.Error())
		RES.Error = "Invalid param: " + err.Error()
		return
	}

	if createReq.Balance, err = util.CoinToRaw(createReq.Balance, 8); err != nil {
		logger.ErrorLog(" - CreateHandler util.CoinToRaw err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = s.cbRepo.Update(createReq); err != nil {
		logger.ErrorLog(" - CreateHandler s.cbRepo.Update err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateCreateReq(createReq domain.ColdBalance) error {
	if createReq.CurrencyId == 0 {
		return errors.New("Currency Id")
	}
	if createReq.Name == "" {
		return errors.New("Name")
	}
	if createReq.Type == "" {
		return errors.New("Type")
	}
	if createReq.Address == "" {
		return errors.New("Address")
	}

	if isFireblocksCold(createReq.Type) && createReq.FireblocksName == "" {
		return errors.New("Fireblocks Name")
	}

	return nil
}
