package eth

import(
    "errors"

    "github.com/btcid/wallet-services-backend/pkg/lib/util"

)

type GetBlockCountRes struct {
    Blocks string
}

func generateRpcReq() util.RpcReq {
    return util.RpcReq{
        RpcUser : "testuser",
        Hash    : "hash",
        Arg1    : "arg1",
        Arg2    : "arg2",
        Arg3    : "arg3",
        Nonce   : "nonce",
    }
}

func GetBlockCount() (*GetBlockCountRes, error) {
    res := GetBlockCountRes{ Blocks: "0" }

    rpcReq := generateRpcReq()
    xmlrpc := util.NewXmlRpc("35.187.234.25", "3000", "/rpc")
    err := xmlrpc.XmlRpcCall("EthRpc.GetBlockCount", &rpcReq, &res)

    if err != nil { 
        return &res, err

    } else if res.Blocks == "0" {
        return &res, errors.New("Unexpected error occured in Node.")

    } else {
        return &res, nil
    }
}