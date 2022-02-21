package coldwallet

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type ListRes struct {
	ColdWallets []cb.ColdBalance `json:"cold_wallets"`
	Error       *errs.Error      `json:"error"`
}

func (re *Rest) ListHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES ListRes
		err error
		ctx = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	service := re.svc.ColdWallet
	if RES.ColdWallets, err = service.ListColdWallet(ctx, page, limit); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
