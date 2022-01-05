package rpcrequest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- rpcrequest.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Request successfully updated"

			config.LoadRpcRequestByRpcMethodId(s.rrqRepo, rpcRequest.RpcMethodId)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcrequest.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcRequest); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if err = validateUpdateReq(rpcRequest); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = s.rrqRepo.Update(rpcRequest); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedUpdateRPCRequest.Title})
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
