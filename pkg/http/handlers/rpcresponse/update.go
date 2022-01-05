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
	if rpcResponse.FieldName == "" {
		return errors.New("Field Name")
	}
	if rpcResponse.XMLPath == "" {
		return errors.New("XML Path")
	}
	if rpcResponse.DataTypeTag == "" {
		return errors.New("Data Type Tag")
	}
	if rpcResponse.RpcMethodId == 0 {
		return errors.New("RPC Method Id")
	}

	return nil
}
