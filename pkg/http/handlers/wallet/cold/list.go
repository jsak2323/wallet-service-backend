package cold

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type ListRes struct {
	ColdWallets []cb.ColdBalance `json:"cold_wallets"`
	Error       *errs.Error      `json:"error"`
}

func (s *ColdWalletService) ListHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES ListRes
		err error
		ctx = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	if RES.ColdWallets, err = s.cbRepo.GetAll(ctx, page, limit); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllColdBalance)
		return
	}
}
