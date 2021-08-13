package ethxmlrpc

import(
    "errors"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
    "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func (es *EthService) GetBalance(rpcConfig rc.RpcConfig) (*model.GetBalanceRpcRes, error) {
    balanceRes := struct {Value string}{}
    res := model.GetBalanceRpcRes{ Balance: "0" }

    rpcReq := util.GenerateRpcReq(rpcConfig, "", "", "")
    xmlrpc := util.NewXmlRpcClient(rpcConfig.Host, rpcConfig.Port, rpcConfig.Path)

    err := xmlrpc.XmlRpcCall("balance", &rpcReq, &balanceRes)

    res.Balance = balanceRes.Value

    if err == nil {
        return &res, nil

    } else if err != nil {
        return &res, err

    } else {
        return &res, errors.New("Unexpected error occured in Node.")
    }
}