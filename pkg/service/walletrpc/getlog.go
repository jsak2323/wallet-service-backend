package walletrpc

import (
	"context"
	"net/http"
	"strings"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *walletRpcService) GetLog(ctx context.Context, SYMBOL string, TOKENTYPE string, rpcConfigType string, date string) (resp *http.Response, err error) {

	SYMBOL = strings.ToUpper(SYMBOL)
	TOKENTYPE = strings.ToUpper(TOKENTYPE)

	currencyConfig, err := config.GetCurrencyBySymbolTokenType(SYMBOL, TOKENTYPE)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetCurrencyBySymbolTokenType)
		return nil, err
	}

	rpcConfig, err := config.GetRpcConfigByType(currencyConfig.Id, rpcConfigType)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetRpcConfigByType)
		return nil, err
	}

	res, err := http.Get("http://" + rpcConfig.Host + ":" + rpcConfig.Port + "/log/" + date)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetLogFile)
		return nil, err
	}

	defer res.Body.Close()

	return res, nil
}
