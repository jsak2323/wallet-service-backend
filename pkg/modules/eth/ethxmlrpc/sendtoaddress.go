package ethxmlrpc

import (
    "errors"
    "encoding/hex"

    "github.com/btcid/wallet-services-backend-go/cmd/config"
    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (es *EthService) SendToAddress(rpcConfig rc.RpcConfig, amountInDecimal string, address string, memo string) (*model.SendToAddressRpcRes, error) {
    txRes := struct {Value string}{}
    res := model.SendToAddressRpcRes{}

    encryptedEncryptKey, err := hex.DecodeString(config.CONF.EthEncryptKeyEncrypted)
    if err != nil { return &res, err }

    decryptedencryptedEncryptKey, err := util.Decrypt(encryptedEncryptKey, []byte(config.CONF.EthEncryptKeyKey))
    if err != nil { return &res, err }
    
    rpcReq := util.GenerateRpcReq(rpcConfig, string(decryptedencryptedEncryptKey), address, amountInDecimal)
    xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err = xmlrpc.XmlRpcCall("send_transaction", &rpcReq, &txRes)

    res.TxHash = txRes.Value

    if err == nil {
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}