package rpcrequest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcRequestService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
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
			logger.InfoLog(" -- rpcrequest.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Request successfully updated"

			config.LoadRpcRequestByRpcMethodId(s.rrqRepo, rpcRequest.RpcMethodId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcrequest.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcRequest); err != nil {
		logger.ErrorLog(" -- rpcrequest.UpdateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateUpdateReq(rpcRequest); err != nil {
		logger.ErrorLog(" -- rpcrequest.UpdateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.rrqRepo.Update(rpcRequest); err != nil {
		logger.ErrorLog(" -- rpcrequest.UpdateHandler rrqRepo.Update Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateUpdateReq(rpcRequest domain.RpcRequest) error {
	if rpcRequest.Id == 0 {
		return errors.New("ID")
	}
	if rpcRequest.Type != domain.TypeJsonRoot && rpcRequest.ArgName == "" {
		return errors.New("Arg Name")
	}
	if rpcRequest.Type == "" {
		return errors.New("Type")
	}
	if rpcRequest.Type != domain.TypeJsonRoot && rpcRequest.Source == "" {
		return errors.New("Source")
	}
	if rpcRequest.RpcMethodId == 0 {
		return errors.New("RPC Method Id")
	}

	return nil
}
