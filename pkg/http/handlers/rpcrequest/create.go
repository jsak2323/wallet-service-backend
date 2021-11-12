package rpcrequest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcRequestService) CreateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcRequest domain.RpcRequest
		RES        StandardRes
		err        error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- rpcrequest.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc Request successfully created"

			config.LoadRpcRequestByRpcMethodId(s.rrqRepo, rpcRequest.RpcMethodId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcrequest.CreateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcRequest); err != nil {
		logger.ErrorLog(" -- rpcrequest.CreateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateCreateReq(rpcRequest); err != nil {
		logger.ErrorLog(" -- rpcrequest.CreateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.rrqRepo.Create(rpcRequest); err != nil {
		logger.ErrorLog(" -- rpcrequest.CreateHandler rrqRepo.Create Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateCreateReq(rpcRequest domain.RpcRequest) error {
	if rpcRequest.ArgName == "" {
		return errors.New("Arg Name")
	}
	if rpcRequest.Type == "" {
		return errors.New("Type")
	}
	if rpcRequest.Source == "" {
		return errors.New("Source")
	}
	if rpcRequest.RpcMethodId == 0 {
		return errors.New("RPC Method Id")
	}

	return nil
}
