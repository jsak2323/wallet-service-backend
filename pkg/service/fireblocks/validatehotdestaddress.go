package fireblocks

import (
	"context"
	"errors"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers/fireblocks"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *fireblocksService) ValidateHotDestAddress(ctx context.Context, req fireblocks.FireblocksSignReq) (resp fireblocks.FireblocksSignRes, err error) {

	resp.Action = fireblocks.ApproveTransaction

	if err = s.validator.Validate(req); err != nil {
		resp.Action = fireblocks.RejectTransaction
		resp.RejectionReason = err.Error()
		err = errs.AssignErr(errs.AddTrace(err), errs.InvalidRequest)
		return resp, err
	}

	coldBalance, err := s.coldbalance.GetByFireblocksName(ctx, req.Asset)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InternalServerErr)
		return resp, err
	}

	currencyRPC := config.CURRRPC[coldBalance.CurrencyId]

	receiverWallet, err := config.GetRpcConfigByType(currencyRPC.Config.Id, rc.SenderRpcType)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.InternalServerErr)
		return resp, err
	}

	if receiverWallet.Address != req.DestAddress {
		resp.Action = fireblocks.RejectTransaction
		resp.RejectionReason = fireblocks.InvalidDestAddressReason
		err = errs.AssignErr(errs.AddTrace(errors.New(fireblocks.InvalidDestAddressReason)), errs.InternalServerErr)
		return resp, err
	}

	return resp, nil

}
