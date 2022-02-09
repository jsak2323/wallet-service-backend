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

func (re *Rest) ListHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES handlerRpcConfig.ListRes
		err error
		ctx = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" - rpcconfig.ListHandler, success!", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	page, _ := strconv.Atoi(vars["page"])
	limit, _ := strconv.Atoi(vars["limit"])

	logger.InfoLog(" - rpcconfig.ListHandler, Requesting ...", req)

	service := re.svc.RpcConfig
	if RES.RpcConfigs, err = service.List(ctx, page, limit); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
