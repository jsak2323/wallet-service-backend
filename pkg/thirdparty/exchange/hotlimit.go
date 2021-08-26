package exchange

import (
	"errors"
	"strings"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type exchangeHotLimitRepository struct {}

func NewExchangeHotLimitRepository() domain.Repository {
	return &exchangeHotLimitRepository{}
}

func (r *exchangeHotLimitRepository) GetBySymbol(symbol string) (result domain.HotLimit, err error) {
	symbol = strings.ToLower(symbol)
	result = make(domain.HotLimit, 5)
	
	if symbol == "btc" {
		result["bottom_hard"] = "10000000000";
		result["bottom_soft"] = "20000000000";
		result["target"] = "30000000000";
		result["top_soft"] = "40000000000";
		result["top_hard"] = "50000000000";

		return result, nil
	}

	if symbol == "usdt" {
		result["bottom_hard"] = "10000000000";
		result["bottom_soft"] = "15000000000";
		result["target"] = "20000000000";
		result["top_soft"] = "25000000000";
		result["top_hard"] = "30000000000";

		return result, nil
	}
	
	sumRes, err := summaries()
	if err != nil {
		return domain.HotLimit{}, err
	}

	ticker, ok := sumRes.Tickers[symbol+"_idr"]
	if !ok {
		return domain.HotLimit{}, errors.New("market trade not found")
	}

	volume := ticker.VolIdr

	if cmp, err := util.CmpBig(volume, "500000000"); err != nil {
		return domain.HotLimit{}, err
	} else if cmp == -1 { // <500jt
		result["bottom_hard"] = "250000000"; //bottom hard Limit
		result["bottom_soft"] = "350000000"; //bottom soft Limit
		result["target"] = "500000000"; //target
		result["top_soft"] = "750000000"; //top soft Limit
		result["top_hard"] = "1000000000"; //top hard limit

		return result, nil
	}

	if cmp, err := util.CmpBig(volume, "1000000000"); err != nil {
		return domain.HotLimit{}, err
	} else if cmp == -1 { // <1m
		result["bottom_hard"] = "500000000"; //bottom hard Limit
		result["bottom_soft"] = "700000000"; //bottom soft Limit
		result["target"] = "1000000000"; //target
		result["top_soft"] = "1500000000"; //top soft Limit
		result["top_hard"] = "2000000000"; //top hard limit

		return result, nil
	}

	if cmp, err := util.CmpBig(volume, "2500000000"); err != nil {
		return domain.HotLimit{}, err
	} else if cmp == -1 {  // <2.5m
		result["bottom_hard"] = "1000000000";
		result["bottom_soft"] = "1500000000";
		result["target"] = "2500000000";
		result["top_soft"] = "4000000000";
		result["top_hard"] = "5000000000";

		return result, nil
	}

	if cmp, err := util.CmpBig(volume, "5000000000"); err != nil {
		return domain.HotLimit{}, err
	} else if cmp == -1 { // <5m
		result["bottom_hard"] = "2500000000";
		result["bottom_soft"] = "3500000000";
		result["target"] = "5000000000";
		result["top_soft"] = "7000000000";
		result["top_hard"] = "10000000000";

		return result, nil
	}

	// >5m
	result["bottom_hard"] = "5000000000";
	result["bottom_soft"] = "7000000000";
	result["target"] = "10000000000";
	result["top_soft"] = "15000000000";
	result["top_hard"] = "20000000000";

	return result, nil
}
