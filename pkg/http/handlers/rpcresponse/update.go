package rpcresponse

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcResponseService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
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
			logger.InfoLog(" -- rpcresponse.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Response successfully updated"

			config.LoadRpcResponseByRpcMethodId(s.rrsRepo, rpcResponse.RpcMethodId)
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcresponse.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcResponse); err != nil {
		logger.ErrorLog(" -- rpcresponse.UpdateHandler json.NewDecoder err: " + err.Error())
		RES.Error = errInternalServer
		return
	}

	if err = validateUpdateReq(rpcResponse); err != nil {
		logger.ErrorLog(" -- rpcresponse.UpdateHandler invalid request: " + err.Error())
		RES.Error = "Invalid request: " + err.Error()
		return
	}

	if err = s.rrsRepo.Update(rpcResponse); err != nil {
		logger.ErrorLog(" -- rpcresponse.UpdateHandler rrsRepo.Update Error: " + err.Error())
		RES.Error = errInternalServer
		return
	}
}

func validateUpdateReq(rpcResponse domain.RpcResponse) error {
	if rpcResponse.Id == 0 {
		return errors.New("ID")
	}
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
		return errors.New("RPC Method Id")
	}

	return nil
}
