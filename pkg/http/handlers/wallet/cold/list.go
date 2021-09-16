package cold

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type ListRes struct {
	ColdWallets []cb.ColdBalance `json:"cold_wallets"`
	Error       string           `json:"error"`
}

func (s *ColdWalletService) ListHandler(w http.ResponseWriter, req *http.Request) {
	var RES ListRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	if RES.ColdWallets, err = s.cbRepo.GetAll(page, limit); err != nil {
		logger.ErrorLog(" - ListHandler s.roleRepo.GetAll err: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}
