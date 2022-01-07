package rpcresponse

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- rpcresponse.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc Response successfully created"

			config.LoadRpcResponseByRpcMethodId(s.rrsRepo, rpcResponse.RpcMethodId)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcresponse.CreateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = validateCreateReq(rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.rrsRepo.Create(rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCResponse)
		return
	}
}

func validateCreateReq(rpcResponse domain.RpcResponse) error {
	if rpcResponse.TargetFieldName == "" {
		return errs.AddTrace(errors.New("Target Field Name"))
	}
	if rpcResponse.XMLPath == "" {
		return errs.AddTrace(errors.New("XML Path"))
	}
	if rpcResponse.DataTypeXMLTag == "" {
		return errs.AddTrace(errors.New("Data Type XML Tag"))
	}
	if rpcResponse.ParseType == "" {
		return errs.AddTrace(errors.New("Parse Type"))
	}
	if rpcResponse.RpcMethodId == 0 {
		return errs.AddTrace(errors.New("Rpc Method Id"))
	}

	err := json.Unmarshal([]byte(rpcResponse.JsonFieldsStr), rpcResponse.JsonFields)
	if err != nil && rpcResponse.JsonFieldsStr != "" {
		return errs.AddTrace(errors.New("JSON data at JSON Field"))
	}

	return nil
}
