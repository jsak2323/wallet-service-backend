package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcConfigService) GetByIdHandler(w http.ResponseWriter, req *http.Request) {
	var RES GetRes
	var err error
	var reqId int

	vars := mux.Vars(req)
	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" - rpcconfig.GetByIdHandler, success!", req)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - rpcconfig.GetByIdHandler, Requesting ...", req)

	if reqId, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByIdHandler strconv.Atoi(" + vars["id"] + ") Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RES.RpcConfig, err = s.rcRepo.GetById(reqId); err != nil {
		logger.ErrorLog(" -- rpcconfig.GetByIdHandler rcRepo.GetById Error: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
