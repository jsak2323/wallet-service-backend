package rpcmethod

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcMethodService) CreateHandler(w http.ResponseWriter, req *http.Request) {
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
			logger.InfoLog(" -- rpcmethod.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc Method successfully created"

			config.LoadRpcMethodByRpcConfigId(s.rmRepo, rpcMethod.RpcConfigId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcmethod.CreateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcMethod); err != nil {
		logger.ErrorLog(" -- rpcmethod.CreateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateCreateReq(rpcMethod); err != nil {
		logger.ErrorLog(" -- rpcmethod.CreateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if rpcMethod.Id, err = s.rmRepo.Create(rpcMethod); err != nil {
		logger.ErrorLog(" -- rpcmethod.CreateHandler rmRepo.Create Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = s.rcrmRepo.Create(rpcMethod.RpcConfigId, rpcMethod.Id); err != nil {
		logger.ErrorLog(" -- rpcmethod.CreateHandler rcrmRepo.Create Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateCreateReq(rpcMethod domain.RpcMethod) error {
	if rpcMethod.Name == "" {
		return errors.New("Name")
	}
	if rpcMethod.Type == "" {
		return errors.New("Type")
	}
	if rpcMethod.RpcConfigId == 0 {
		return errors.New("RPC Config ID")
	}

	return nil
}
