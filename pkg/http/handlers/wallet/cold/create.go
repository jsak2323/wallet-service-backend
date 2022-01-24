package cold

import (
	"encoding/json"
	"errors"
	"net/http"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (s *ColdWalletService) CreateHandler(w http.ResponseWriter, req *http.Request) {
	var createReq domain.CreateColdBalance
	var RES StandardRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
		} else {
			RES.Success = true
			RES.Message = "Cold wallet successfully created"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	if err = json.NewDecoder(req.Body).Decode(&createReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = s.validator.Validate(createReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = validateCreateReq(createReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	createReq.Balance = util.CoinToRaw(createReq.Balance, 8)

	if _, err = s.cbRepo.Create(createReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateColdBalance)
		return
	}
}

func validateCreateReq(createReq domain.CreateColdBalance) error {

	if isFireblocksCold(createReq.Type) && createReq.FireblocksName == "" {
		return errs.AddTrace(errors.New("Fireblocks Name"))
	}

	return nil
}
