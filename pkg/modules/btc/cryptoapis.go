package btc

import(
    "encoding/json"

    "gopkg.in/resty.v0"

    "github.com/btcid/wallet-services-backend/cmd/config"
)

type CryptoApisGetNodeInfoRes struct {
    Payload CryptoApisGetNodeInfoResPayload `json:"payload"`
}

type CryptoApisGetNodeInfoResPayload struct {
    Difficulty              float64     `json:"difficulty"`
    Headers                 int         `json:"headers"`
    Chain                   string      `json:"chain"`
    Chainwork               string      `json:"chainwork"`
    MedianTime              int         `json:"mediantime"`
    Blocks                  int         `json:"blocks"`
    BestBlockHash           string      `json:"bestblockhash"`
    Currency                string      `json:"currency"`
    Transactions            int         `json:"transactions"`
    VerificationProgress    float64     `json:"verificationprogress"`
}

type CryptoApisService struct{
    ApiKey      string
    Network     string
}

func NewCryptoApisService() *CryptoApisService {
    network := "mainnet"
    if config.IS_DEV {
        network = "testnet"
    }

    return &CryptoApisService{
        ApiKey      : config.CONF.CryptoApisKey,
        Network     : network,
    }
}

func (cas *CryptoApisService) GetNodeInfo() (CryptoApisGetNodeInfoRes, error) {
    getNodeInfoRes := CryptoApisGetNodeInfoRes{}

    restRes, err := resty.R().
      SetHeader("Content-Type", "application/json").
      SetHeader("X-API-Key", cas.ApiKey).
      Get("https://api.cryptoapis.io/v1/bc/btc/"+cas.Network+"/info")

    if err == nil {
        err = json.Unmarshal([]byte(restRes.String()), &getNodeInfoRes)
    }

    return getNodeInfoRes, err
}


