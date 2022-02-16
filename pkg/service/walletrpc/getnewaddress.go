package walletrpc

import (
	"context"
	"strings"

	config "github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *walletRpcService) GetNewAddress(ctx context.Context, symbol string, tokenType string, addressType string) (resp *handlers.GetNewAddressRes, err error) {
	SYMBOL := strings.ToUpper(symbol)

	currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, tokenType)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
		return nil, err
	}

	rpcConfig, err := config.GetRpcConfigByType(currencyConfig.Id, "receiver")
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRpcConfigByType)
		return nil, err
	}
	resp.RpcConfig = handlers.RpcConfigResDetail{
		RpcConfigId:          rpcConfig.Id,
		Symbol:               SYMBOL,
		Name:                 rpcConfig.Name,
		Host:                 rpcConfig.Host,
		Type:                 rpcConfig.Type,
		NodeVersion:          rpcConfig.NodeVersion,
		NodeLastUpdated:      rpcConfig.NodeLastUpdated,
		IsHealthCheckEnabled: rpcConfig.IsHealthCheckEnabled,
	}

	module, err := s.moduleServices.GetModule(currencyConfig.Id)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetModule)
		return nil, err
	}

	rpcRes, err := module.GetNewAddress(ctx, rpcConfig, addressType)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetNewAddress)
		return nil, err
	}

	resp.Address = rpcRes.Address
	return resp, err
}
