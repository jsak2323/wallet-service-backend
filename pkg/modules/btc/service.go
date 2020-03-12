package btc

import(
    "errors"

    m "github.com/btcid/wallet-services-backend/pkg/modules"
)

func GetBlockCount(rpcConfig rc.RpcConfig) (*m.GetBlockCountRpcRes, error) {
    res := m.GetBlockCountRpcRes{ Blocks: "0" }

    rpcReq := m.GenerateRpcReq(rpcConfig)
    xmlrpc := util.NewXmlRpc(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)
    
    err := xmlrpc.XmlRpcCall("bitrpc.getblockcount", &rpcReq, &res)

    if err != nil { 
        return &res, err

    } else if res.Blocks == "0" {
        return &res, errors.New("Unexpected error occured in Node.")

    } else {
        return &res, nil
    }
}