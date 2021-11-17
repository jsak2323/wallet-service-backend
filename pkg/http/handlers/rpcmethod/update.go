package rpcmethod

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcMethodService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcMethod domain.RpcMethod
		RES       StandardRes
		err       error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- rpcmethod.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Method successfully updated"

			config.LoadRpcMethodByRpcConfigId(s.rmRepo, rpcMethod.RpcConfigId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcmethod.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcMethod); err != nil {
		logger.ErrorLog(" -- rpcmethod.UpdateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateUpdateReq(rpcMethod); err != nil {
		logger.ErrorLog(" -- rpcmethod.UpdateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.rmRepo.Update(rpcMethod); err != nil {
		logger.ErrorLog(" -- rpcmethod.UpdateHandler rmRepo.Update Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateUpdateReq(rpcMethod domain.RpcMethod) error {
	if rpcMethod.Id == 0 {
		return errors.New("ID")
	}
	if rpcMethod.Name == "" {
		return errors.New("Name")
	}
	if rpcMethod.Type == "" {
		return errors.New("Type")
	}

	return nil
}
