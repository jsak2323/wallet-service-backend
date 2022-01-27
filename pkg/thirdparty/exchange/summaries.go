package exchange

import (
	"encoding/json"

	"gopkg.in/resty.v0"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

type SummariesRes struct {
	Tickers   map[string]Ticker `json:"tickers"`
	Prices24H map[string]string `json:"prices_24h"`
}

type Ticker struct {
	Last   string `json:"last"`
	VolIdr string `json:"vol_idr"`
	Name   string `json:"name"`
}

var host string = config.CONF.ExchangeHost

func summaries() (RES SummariesRes, err error) {
	res, err := resty.R().Get(host + "/api/summaries")
	if err != nil {
		return SummariesRes{}, errs.AddTrace(err)
	}

	if err = json.Unmarshal(res.Body, &RES); err != nil {
		return SummariesRes{}, errs.AddTrace(err)
	}

	return RES, nil
}
