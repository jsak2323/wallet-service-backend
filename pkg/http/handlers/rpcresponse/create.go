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
		if err != nil {
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
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if err = validateCreateReq(rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = s.rrsRepo.Create(rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedCreateRPCResponse.Title})
		return
	}
}

func validateCreateReq(rpcResponse domain.RpcResponse) error {
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
		return errors.New("Rpc Method Id")
	}

	return nil
}
