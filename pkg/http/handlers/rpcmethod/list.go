package rpcmethod

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcMethodService) ListHandler(w http.ResponseWriter, req *http.Request) {
	var RES ListRes
	var err error

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" - rpcmethod.ListHandler, success!", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	logger.InfoLog(" - rpcmethod.ListHandler, Requesting ...", req)

	if RES.RpcMethods, err = s.rmRepo.GetAll(page, limit); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllRPCMethod)
		return
	}
}
