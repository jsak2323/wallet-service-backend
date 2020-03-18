package modules

import(
    modules_m "github.com/btcid/wallet-services-backend/pkg/modules/model"
    "github.com/btcid/wallet-services-backend/pkg/modules/btc"
    "github.com/btcid/wallet-services-backend/pkg/modules/eth"
)

var (
    MS = make(map[string]modules_m.ModuleService)
    ModuleServices = &MS
)

func init() {
    (*ModuleServices)["BTC"] = &btc.BtcService{}
    (*ModuleServices)["ETH"] = &eth.EthService{}
}