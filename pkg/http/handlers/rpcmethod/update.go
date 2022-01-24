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
		rpcMethod domain.UpdateRpcMethod
		RES       StandardRes
		err       error
		ctx       = req.Context()
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
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
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = s.validator.Validate(rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if err = s.rmRepo.Update(rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateRPCMethod)
		return
	}
}

func validateUpdateReq(rpcMethod domain.RpcMethod) error {
	if rpcMethod.Id == 0 {
		return errs.AddTrace(errors.New("ID"))
	}
	if rpcMethod.Name == "" {
		return errs.AddTrace(errors.New("Name"))
	}
	if rpcMethod.Type == "" {
		return errs.AddTrace(errors.New("Type"))
	}

	return nil
}
