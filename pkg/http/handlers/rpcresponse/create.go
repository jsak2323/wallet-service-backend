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
		rpcResponse domain.CreateRpcResponse
		RES         StandardRes
		err         error
		ctx         = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- rpcresponse.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc Response successfully created"

			config.LoadRpcResponseByRpcMethodId(ctx, s.rrsRepo, rpcResponse.RpcMethodId)
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
	if err = s.validator.Validate(rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.rrsRepo.Create(ctx, rpcResponse); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCResponse)
		return
	}
}

func validateCreateReq(rpcResponse domain.CreateRpcResponse) error {

	err := json.Unmarshal([]byte(rpcResponse.JsonFieldsStr), rpcResponse.JsonFields)
	if err != nil && rpcResponse.JsonFieldsStr != "" {
		return errs.AddTrace(errors.New("JSON data at JSON Field"))
	}

	return nil
}
