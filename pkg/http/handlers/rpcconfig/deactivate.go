package rpcconfig

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/gorilla/mux"
)

func (svc *RpcConfigService) DeactivateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		id  int
		RES StandardRes
		err error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- rpcconfig.DeactivateHandler Success!", req)

			RES.Success = true
			RES.Message = "RpcConfig successfully deactivated"

			config.LoadAppConfig()
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	vars := mux.Vars(req)
	if id, err = strconv.Atoi(vars["id"]); err != nil {
		logger.ErrorLog(" -- rpcconfig.DeactivateHandler invalid request")
		RES.Error = "Invalid request"
	}

	if err = svc.rcRepo.ToggleActive(id, false); err != nil {
		logger.ErrorLog(" -- rpcconfig.DeactivateHandler svc.ccRepo.ToggleActive err: " + err.Error())
		RES.Error = err.Error()
		return
	}
}
