package rpcconfig

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	handlerRpcConfig "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (re *Rest) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcConfig domain.UpdateRpcConfig
		RES       handlerRpcConfig.StandardRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- rpcconfig.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "Currency successfully updated"

			config.LoadCurrencyConfigs(ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcconfig.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcConfig); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.RpcConfig
	if err = service.Update(ctx, rpcConfig); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}
}
