package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcConfigService) ListHandler(w http.ResponseWriter, req *http.Request) {
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

	if RES.RpcConfigs, err = s.rcRepo.GetAll(ctx, page, limit); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetAllRPCConfig)
		return
	}
}
