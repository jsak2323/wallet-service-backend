package xmlrpc

import (
    "errors"
    "strings"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type SendToAddressXmlRpcRes struct {
    Content SendToAddressXmlRpcResStruct
}
type SendToAddressXmlRpcResStruct struct {
    TxHash  string
    Error   string
}

func (gts *GeneralTokenService) SendToAddress(rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*model.SendToAddressRpcRes, error) {
    res := model.SendToAddressRpcRes{}

    token := strings.ToLower(gts.Symbol)

    rpcReq := util.GenerateRpcReq(rpcConfig, token, amountInDecimal, address)
    client := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    rpcRes := SendToAddressXmlRpcRes{}

    err := client.XmlRpcCall(gts.ParentSymbol+"Rpc.SendToAddress", &rpcReq, &rpcRes)

    if err != nil {
        return &res, err

    } else if rpcRes.Content.Error != "" { 
        return &res, errors.New(rpcRes.Content.Error)

    } else if rpcRes.Content.TxHash == "" {
        return &res, errors.New("Unexpected error occured in Node.")
    }

    res.TxHash = rpcRes.Content.TxHash
    return &res, nil
}


