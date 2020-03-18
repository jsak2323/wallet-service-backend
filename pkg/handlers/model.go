package handlers

import(
	modules_m "github.com/btcid/wallet-services-backend/pkg/modules/model"
	"github.com/btcid/wallet-services-backend/pkg/modules/btc"
	"github.com/btcid/wallet-services-backend/pkg/modules/eth"
)

type GetBlockCountRes struct{
    Symbol              string
    Host                string
    Type                string
    NodeVersion         string
    NodeLastUpdated     string
    Blocks              string
}

func NewModuleServices() *map[string]modules_m.ModuleService {
    ModuleServices := make(map[string]modules_m.ModuleService)

    ModuleServices["BTC"] = &btc.BtcService{}
    ModuleServices["ETH"] = &eth.EthService{}

    return &ModuleServices
}