package currency

import (
	"errors"
)

func (req CurrencyRpcReq) validate() error {
	if req.CurrencyId == 0 {
		return errors.New("Invalid Currency ID")
	}
	
	if req.RpcId == 0 {
		return errors.New("Invalid Rpc Config ID")
	}
	
	return nil
}