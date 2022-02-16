package walletrpc

import (
	"context"
	"strings"

	config "github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type SendToAddressRequest struct {
	Symbol    string `json:"symbol"`
	TokenType string `json:"token_type"`
	Amount    string `json:"amount"`
	Address   string `json:"address"`
	Memo      string `json:"memo"`
}

func (s *walletRpcService) SendToAddress(ctx context.Context, req handlers.SendToAddressRequest) (resp *handlers.SendToAddressRes, err error) {

	var (
		amountInDecimal = req.Amount
		SYMBOL          = strings.ToUpper(req.Symbol)
		TOKENTYPE       = strings.ToUpper(req.TokenType)
	)

	currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, req.TokenType)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
		return resp, err
	}

	// define rpc config
	rpcConfig, err := config.GetRpcConfigByType(currencyConfig.Id, "sender")
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRpcConfigByType)
		return resp, err
	}
	resp.RpcConfig = handlers.RpcConfigResDetail{
		RpcConfigId:          rpcConfig.Id,
		Symbol:               SYMBOL,
		TokenType:            TOKENTYPE,
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

	// execute rpc call
	rpcRes, err := module.SendToAddress(ctx, rpcConfig, amountInDecimal, req.Address, req.Memo)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedSendToAddress)
		return resp, err
	}

	resp.TxHash = rpcRes.TxHash
	return resp, nil

}
