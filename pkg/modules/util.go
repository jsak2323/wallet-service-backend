package modules

import(
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

func GenerateRpcReq(rpcConfig rc.RpcConfig, arg1 string, arg2 string, arg3 string) util.RpcReq {
    
    return util.RpcReq{
        RpcUser : rpcConfig.User,
        Hash    : rpcConfig.Hashkey,
        Arg1    : "arg1",
        Arg2    : "arg2",
        Arg3    : "arg3",
        Nonce   : "nonce",
    }
}

func NewModuleServices() *map[string]ModuleService {
    ModuleServices := make(map[string]ModuleService)

    ModuleServices["BTC"] = &BtcService{}
    ModuleServices["ETH"] = &EthService{}

    return &ModuleServices
}