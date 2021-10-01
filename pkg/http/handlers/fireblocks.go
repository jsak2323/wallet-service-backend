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

	RES.Action = ApproveTransaction

	if err = json.NewDecoder(req.Body).Decode(&SignReq); err != nil {
		logger.ErrorLog(" -- fireblocks.CallbackHandler json.NewDecoder err: " + err.Error())
		RES.Action = RejectTransaction
		RES.RejectionReason = errInternalServer
		return
	}

	if err := SignReq.Validate(); err != nil {
		RES.Action = RejectTransaction
		RES.RejectionReason = "Invalid param: " + err.Error()
		return
	}

	if SignReq.DestId == config.CONF.FireblocksHotVaultId {
		validateHotDestAddress(SignReq, &RES)
	}
}

func validateHotDestAddress(signReq FireblocksSignReq, res *FireblocksSignRes) {
	receiverWallet, err := rc.GetReceiverFromList(config.CURR[signReq.Asset].RpcConfigs)
	if err != nil {
		logger.ErrorLog(" -- fireblocks.CallbackHandler rc.GetReceiverFromList err: " + err.Error())
		res.Action = RejectTransaction
		res.RejectionReason = errInternalServer
		return
	}
	
	if receiverWallet.Address != signReq.DestAddress {
		res.Action = RejectTransaction
		res.RejectionReason = InvalidDestAddressReason
	}
}

func (r *FireblocksSignReq) Validate() (err error) {
	if r.Asset == "" {
		return errors.New("asset is required")
	}

	if r.DestId == "" {
		return errors.New("destId is required")
	}

	if r.DestAddress == "" {
		return errors.New("destAddress is required")
	}
	
	return nil
}