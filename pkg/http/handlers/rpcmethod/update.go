package rpcmethod

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *RpcMethodService) UpdateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcMethod domain.RpcMethod
		RES       StandardRes
		err       error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- rpcmethod.UpdateHandler Success!", req)

			RES.Success = true
			RES.Message = "RPC Method successfully updated"

			config.LoadRpcMethodByRpcConfigId(s.rmRepo, rpcMethod.RpcConfigId)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcmethod.UpdateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.ErrorUnmarshalBodyRequest.Title})
		return
	}

	if err = validateUpdateReq(rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.InvalidRequest.Title})
		return
	}

	if err = s.rmRepo.Update(rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), &errs.Error{Title: errs.FailedUpdateRPCMethod.Title})
		return
	}
}

func validateUpdateReq(rpcMethod domain.RpcMethod) error {
	if rpcMethod.Id == 0 {
		return errors.New("ID")
	}
	if rpcMethod.Name == "" {
		return errors.New("Name")
	}
	if rpcMethod.Type == "" {
		return errors.New("Type")
	}

	return nil
}
