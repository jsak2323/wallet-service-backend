package coldwallet

import (
	"context"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	handlers "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
)

func (s *coldWalletService) SendToHot(ctx context.Context, sendToHotReq handlers.SendToHotReq) (err error) {
	vaultAccountId, err := s.FireblocksVaultAccountId(sendToHotReq.FireblocksType)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedFireblocksVaultAccountId)
		return err
	}

	res, err := fireblocks.CreateTransaction(fireblocks.CreateTransactionReq{
		AssetId: sendToHotReq.FireblocksName,
		Amount:  sendToHotReq.Amount,
		Source: fireblocks.TransactionAccount{
			Type: fireblocks.VaultAccountType,
			Id:   vaultAccountId,
		},
		Destination: fireblocks.TransactionAccount{
			Type: fireblocks.InternalWalletType,
			Id:   config.CONF.FireblocksHotVaultId,
		},
	})
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateTransaction)
		return err
	}

	if res.Error != "" {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedCreateTransaction)
		return err
	}

	return nil
}
