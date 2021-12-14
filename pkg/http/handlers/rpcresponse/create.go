package rpcresponse

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcResponseService) CreateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcResponse domain.RpcResponse
		RES         StandardRes
		err         error
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != "" {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- rpcresponse.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc Response successfully created"

			config.LoadRpcResponseByRpcMethodId(s.rrsRepo, rpcResponse.RpcMethodId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcresponse.CreateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcResponse); err != nil {
		logger.ErrorLog(" -- rpcresponse.CreateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateCreateReq(rpcResponse); err != nil {
		logger.ErrorLog(" -- rpcresponse.CreateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.rrsRepo.Create(rpcResponse); err != nil {
		logger.ErrorLog(" -- rpcresponse.CreateHandler rrsRepo.Create Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateCreateReq(rpcResponse domain.RpcResponse) error {
	if rpcResponse.TargetFieldName == "" {
		return errors.New("Target Field Name")
	}
	if rpcResponse.XMLPath == "" {
		return errors.New("XML Path")
	}
	if rpcResponse.DataTypeXMLTag == "" {
		return errors.New("Data Type XML Tag")
	}
	if rpcResponse.ParseType == "" {
		return errors.New("Parse Type")
	}
	if rpcResponse.RpcMethodId == 0 {
		return errors.New("Rpc Method Id")
	}

	return nil
}
