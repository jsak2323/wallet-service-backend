package eth

import (
    "fmt"
    "encoding/json"

    "gopkg.in/resty.v0"

    "github.com/btcid/wallet-services-backend/cmd/config"
)

type InfuraEthBlockNumberRes struct {
    Jsonrpc     string `json:"jsonrpc"` 
    Id          string `json:"id"` 
    Result      string `json:"result"` 
}

type InfuraService struct{
    ProjectId   string
    Network     string
}

func NewInfuraService() *InfuraService {
    network := "mainnet"
    if config.IS_DEV {
        network = "ropsten"
    }

    return &InfuraService{
        ProjectId   : config.CONF.InfuraProjectId,
        Network     : network,
    }
}

func (is *InfuraService) EthBlockNumber() (InfuraEthBlockNumberRes, error) {
    ethBlockNumberRes := InfuraEthBlockNumberRes{}
    
    endpoint := "https://"+is.Network+".infura.io/v3/"+is.ProjectId
    restRes, err := resty.R().
          SetHeader("Content-Type", "application/json").
          SetBody(`{"jsonrpc":"2.0","method":"eth_blockNumber","params": [],"id":1}`).
          Post(endpoint)

    if err == nil {
        err = json.Unmarshal([]byte(restRes.String()), &ethBlockNumberRes)
    }

    return ethBlockNumberRes, err
}