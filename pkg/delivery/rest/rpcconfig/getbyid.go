package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	handlerRpcConfig "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) GetByIdHandler(w http.ResponseWriter, req *http.Request) {
	var RES handlerRpcConfig.GetRes
	var err error
	var reqId int
	var ctx = req.Context()

	vars := mux.Vars(req)
	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" - rpcconfig.GetByIdHandler, success!", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - rpcconfig.GetByIdHandler, Requesting ...", req)

	if reqId, err = strconv.Atoi(vars["id"]); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	service := re.svc.RpcConfig
	if RES.RpcConfig, err = service.GetById(ctx, reqId); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
