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

func (s *RpcResponseService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcResponse domain.RpcResponse
		RES         StandardRes
		err         error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {

			logger.InfoLog(" -- rpcresponse.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Response successfully updated"

			config.LoadRpcResponseByRpcMethodId(s.rrsRepo, rpcResponse.RpcMethodId)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcresponse.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if err = validateUpdateReq(rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = s.rrsRepo.Update(rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedUpdateRPCResponse.Title})
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

	err := json.Unmarshal([]byte(rpcResponse.JsonFieldsStr), rpcResponse.JsonFields)
	if err != nil && rpcResponse.JsonFieldsStr != "" {
		return errors.New("JSON data at JSON Field")
	}

	return nil
}
