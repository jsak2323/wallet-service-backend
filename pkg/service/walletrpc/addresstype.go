package walletrpc

import (
	"context"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	handler "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *walletRpcService) GetAddressType(ctx context.Context, symbol string, tokenType string, address string) (resp handler.AddressTypeRes, err error) {
	SYMBOL := strings.ToUpper(symbol)

	currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, tokenType)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
		return resp, err
	}

	rpcConfig, err := config.GetRpcConfigByType(currencyConfig.Id, "receiver")
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRpcConfigByType)
		return resp, err
	}

	resp.RpcConfig = handler.RpcConfigResDetail{
		RpcConfigId:          rpcConfig.Id,
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
		return resp, err
	}

	rpcRes, err := module.AddressType(ctx, rpcConfig, address)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedAddressType)
		return resp, err
	}

	resp.AddressType = rpcRes.AddressType
	return resp, err
}
