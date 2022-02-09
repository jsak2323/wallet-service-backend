package ethxmlrpc

import (
	"context"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"

	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

func (es *EthService) ListWithdraws(ctx context.Context, rpcConfig rc.RpcConfig, limit int) (*model.ListWithdrawsRpcRes, error) {
	return nil, nil
}
