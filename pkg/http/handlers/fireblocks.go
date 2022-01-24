package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type FireblocksService struct {
}

func NewFireblocksService() *FireblocksService {
	return &FireblocksService{}
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
		if RES.RejectionReason == errInternalServer {
			resStatus = http.StatusInternalServerError
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
		logger.ErrorLog(" -- fireblocks.CallbackHandler json.NewDecoder err: "+err.Error(), ctx)
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
		validateHotDestAddress(ctx, SignReq, &RES)
	}
}

func validateHotDestAddress(ctx context.Context, signReq FireblocksSignReq, res *FireblocksSignRes) {
	var currencyConfig cc.CurrencyConfig

	if signReq.Type == TypeBaseAsset {
		signReq.Type = cc.MainTokenType
	}

	currencyConfig, err := config.GetCurrencyBySymbolTokenType(signReq.Asset, signReq.Type)
	if err != nil {
		logger.ErrorLog(" -- fireblocks.CallbackHandler config.GetCurrencyBySymbol("+signReq.Asset+","+signReq.Type+")+err: "+err.Error(), ctx)
		res.Action = RejectTransaction
		res.RejectionReason = errAssetNotFound
		return
	}

	receiverWallet, err := config.GetRpcConfigByType(currencyConfig.Id, rc.SenderRpcType)
	if err != nil {
		logger.ErrorLog(" -- fireblocks.CallbackHandler rc.GetRpcConfigByType err: "+err.Error(), ctx)
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
		return errs.AddTrace(errors.New("asset is required"))
	}

	if r.Type == "" {
		return errs.AddTrace(errors.New("type is required"))
	}

	if r.DestId == "" {
		return errs.AddTrace(errors.New("destId is required"))
	}

	if r.DestAddress == "" {
		return errs.AddTrace(errors.New("destAddress is required"))
	}

	return nil
}
