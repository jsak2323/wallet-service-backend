package handlers

import (
	"errors"
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type FireblocksService struct {
}

func NewFireblocksService() *FireblocksService {
	return &FireblocksService{}
}

const RejectTransaction = "REJECT"
const ApproveTransaction = "APPROVE"
const IgnoreTransaction = "IGNORE"

const InvalidDestAddressReason = "Invalid destination address"

func (s *FireblocksService) CallbackHandler(w http.ResponseWriter, req *http.Request) {
	var SignReq FireblocksSignReq
	var RES FireblocksSignRes
	var err error

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.RejectionReason == errInternalServer {
			resStatus = http.StatusInternalServerError
		} else {
			logger.InfoLog(" -- fireblocks.CallbackHandler Success!", req)
		}

		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- fireblocks.CallbackHandler, Requesting ...", req)

	if err = json.NewDecoder(req.Body).Decode(&SignReq); err != nil {
		logger.ErrorLog(" -- fireblocks.CallbackHandler json.NewDecoder err: " + err.Error())
		RES.RejectionReason = errInternalServer
		return
	}

	RES.Action = RejectTransaction

	if err = SignReq.Validate(); err != nil {
		logger.ErrorLog(" -- fireblocks.CallbackHandler SignReq.Validate err: " + err.Error())
		RES.RejectionReason = "Invalid param: " + err.Error()
		return
	}

	RES.RejectionReason = InvalidDestAddressReason
	
	receiverWallet, err := rc.GetReceiverFromList(config.CURR[SignReq.Asset].RpcConfigs)
	if err != nil {
		logger.ErrorLog(" -- fireblocks.CallbackHandler rc.GetReceiverFromList err: " + err.Error())
		RES.RejectionReason = errInternalServer
		return
	}
	
	if receiverWallet.Address == SignReq.DestAddress {
		RES.Action = ApproveTransaction
		RES.RejectionReason = ""
	}
}

func (r *FireblocksSignReq) Validate() (err error) {
	if r.Asset == "" {
		return errors.New("Asset is required")
	}

	if r.DestAddress == "" {
		return errors.New("Destination Address is required")
	}
	
	return nil
}