package thetaxmlrpc

import (
    // "errors"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    // "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type SendToAddressXmlRpcRes struct {
    Content SendToAddressXmlRpcResStruct
}
type SendToAddressXmlRpcResStruct struct {
    TxHash  string
    Error   string
}

// to be implemented
func (ts *ThetaService) SendToAddress(rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*model.SendToAddressRpcRes, error) {
    res := model.SendToAddressRpcRes{}
    return &res, nil
}


