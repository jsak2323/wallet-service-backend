package rpcmethod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcMethodService) ListHandler(w http.ResponseWriter, req *http.Request) {
	var RES ListRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" - rpcmethod.ListHandler, success!", req)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	logger.InfoLog(" - rpcmethod.ListHandler, Requesting ...", req)

	if RES.RpcMethods, err = s.rmRepo.GetAll(page, limit); err != nil {
		logger.ErrorLog(" -- rpcmethod.ListHandler rmRepo.GetAll Error: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
