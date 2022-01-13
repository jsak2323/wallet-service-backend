package currency

import (
	"errors"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (req CurrencyRpcReq) validate() error {
	if req.CurrencyId == 0 {
		return errs.AddTrace(errors.New("Invalid Currency ID"))
	}

	if req.RpcId == 0 {
		return errs.AddTrace(errors.New("Invalid Rpc Config ID"))
	}

	return nil
}
