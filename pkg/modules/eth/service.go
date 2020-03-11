package eth

import(
    "errors"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend/cmd/config"
    "github.com/btcid/wallet-services-backend/pkg/lib/util"
)

type GetBlockCountRpcRes struct {
    Blocks string
}

func generateRpcReq(rpcConfig *rc.RpcConfig) util.RpcReq {

    return util.RpcReq{
        RpcUser : rpcConfig.User,
        Hash    : rpcConfig.Hashkey,
        Arg1    : "arg1",
        Arg2    : "arg2",
        Arg3    : "arg3",
        Nonce   : "nonce",
    }

    // return util.RpcReq{
    //     RpcUser : "ifan",
    //     Hash    : "testhashkey",
    //     Arg1    : "arg1",
    //     Arg2    : "arg2",
    //     Arg3    : "arg3",
    //     Nonce   : "nonce",
    // }
}

func GetBlockCount(rpcConfig *rc.RpcConfig) (*GetBlockCountRpcRes, error) {
    res := GetBlockCountRpcRes{ Blocks: "0" }

    // var senderRpcConfig rc.RpcConfig
    // for _, rpcConfig := range config.CURR["ETH"].RpcConfigs {
    //     if rpcConfig.Type == "sender" {
    //         senderRpcConfig = rpcConfig
    //         break
    //     }
    // }    

    rpcReq := generateRpcReq(rpcConfig)
    // xmlrpc := util.NewXmlRpc("35.187.234.25", "3000", "/rpc")
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)
    err := xmlrpc.XmlRpcCall("EthRpc.GetBlockCount", &rpcReq, &res)

    if err != nil { 
        return &res, err

    } else if res.Blocks == "0" {
        return &res, errors.New("Unexpected error occured in Node.")

    } else {
        return &res, nil
    }
}