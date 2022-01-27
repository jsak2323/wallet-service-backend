package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type FireblocksService struct {
	cbRepo cb.Repository
}

func NewFireblocksService(cbRepo cb.Repository) *FireblocksService {
	return &FireblocksService{cbRepo: cbRepo}
}

const TypeBaseAsset = "BASE_ASSET"

const RejectTransaction = "REJECT"
const ApproveTransaction = "APPROVE"
const IgnoreTransaction = "IGNORE"

const InvalidDestAddressReason = "Invalid destination address"

func (s *FireblocksService) CallbackHandler(w http.ResponseWriter, req *http.Request) {
	var (
		SignReq FireblocksSignReq
		RES     FireblocksSignRes
		err     error
		ctx     = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError

			RES.Action = RejectTransaction
			RES.RejectionReason = RES.Error.Error()

			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- fireblocks.CallbackHandler Success!", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- fireblocks.CallbackHandler, Requesting ...", req)

	RES.Action = ApproveTransaction

	if err = json.NewDecoder(req.Body).Decode(&SignReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	if err := SignReq.Validate(); err != nil {
		RES.Action = RejectTransaction
		RES.RejectionReason = "Invalid param: " + err.Error()
		return
	}

	if SignReq.DestId == config.CONF.FireblocksHotVaultId {
		s.validateHotDestAddress(SignReq, &RES)
	}
}

func (s *FireblocksService) validateHotDestAddress(signReq FireblocksSignReq, res *FireblocksSignRes) {
	coldBalance, err := s.cbRepo.GetByFireblocksName(signReq.Asset)
	if err != nil {
		res.Error = errs.AssignErr(errs.AddTrace(err), errs.InternalServerErr)
		return
	}

	currencyRPC := config.CURRRPC[coldBalance.CurrencyId]

	receiverWallet, err := config.GetRpcConfigByType(currencyRPC.Config.Id, rc.SenderRpcType)
	if err != nil {
		res.Error = errs.AssignErr(errs.AddTrace(err), errs.InternalServerErr)
		return
	}

	if receiverWallet.Address != signReq.DestAddress {
		res.Action = RejectTransaction
		res.RejectionReason = InvalidDestAddressReason
	}
}

func (r *FireblocksSignReq) Validate() (err error) {
	if r.Asset == "" {
		return errs.AddTrace(errors.New("asset is required"))
	}

	if r.DestId == "" {
		return errs.AddTrace(errors.New("destId is required"))
	}

	if r.DestAddress == "" {
		return errs.AddTrace(errors.New("destAddress is required"))
	}

	return nil
}
