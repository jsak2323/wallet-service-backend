package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcConfigService) ListHandler(w http.ResponseWriter, req *http.Request) {
	var RES ListRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" - rpcconfig.ListHandler, success!", req)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
    page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	logger.InfoLog(" - rpcconfig.ListHandler, Requesting ...", req)

	if RES.RpcConfigs, err = s.rcRepo.GetAll(page, limit); err != nil {
		logger.ErrorLog(" -- rpcconfig.ListHandler rcRepo.GetAll Error: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
