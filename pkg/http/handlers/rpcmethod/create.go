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

func (s *RpcMethodService) CreateHandler(w http.ResponseWriter, req *http.Request) {
	var (
		rpcMethod domain.RpcMethod
		RES       StandardRes
		err       error
	)

	handleResponse := func() {

		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error))
		} else {
			logger.InfoLog(" -- rpcmethod.CreateHandler Success!", req)

			RES.Success = true
			RES.Message = "Rpc Method successfully created"

			config.LoadRpcMethodByRpcConfigId(s.rmRepo, rpcMethod.RpcConfigId)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- rpcmethod.CreateHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err = validateCreateReq(rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return
	}

	if rpcMethod.Id, err = s.rmRepo.Create(rpcMethod); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCMethod)
		return
	}

	if err = s.rcrmRepo.Create(rpcMethod.RpcConfigId, rpcMethod.Id); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateRPCConfigRPCMethod)
		return
	}
}

func validateCreateReq(rpcMethod domain.RpcMethod) error {
	if rpcMethod.Name == "" {
		return errs.AddTrace(errors.New("Name"))
	}
	if rpcMethod.Type == "" {
		return errs.AddTrace(errors.New("Type"))
	}
	if rpcMethod.RpcConfigId == 0 {
		return errs.AddTrace(errors.New("RPC Config ID"))
	}

	return nil
}
