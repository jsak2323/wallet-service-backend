package btc

import(
    "gopkg.in/resty.v0"

    "github.com/btcid/wallet-services-backend/cmd/config"
)

var CryptoApisApiKey string = config.CONF.CryptoApisKey

type CryptoApisGetNodeInfoRes struct {
    Payload CryptoApisGetNodeInfoResPayload `json:"payload"`
}

type CryptoApisGetNodeInfoResPayload struct {
    Difficulty              string `json:"difficulty"`
    Headers                 string `json:"headers"`
    Chain                   string `json:"chain"`
    Chainwork               string `json:"chainwork"`
    MedianTime              string `json:"mediantime"`
    Blocks                  string `json:"blocks"`
    BestBlockHash           string `json:"bestblockhash"`
    Currency                string `json:"currency"`
    Transactions            string `json:"transactions"`
    VerificationProgress    string `json:"verificationprogress"`
}

type CryptoApisService struct{
    ApiKey      string
    Network     string
}

func NewCryptoApisService() *CryptoApisService {
    network := "mainnet"
    if os.GetEnv("IS_DEV") {
        network := "testnet"
    }

    return &CryptoApisService{
        ApiKey      : CryptoApisApiKey,
        Network     : network,
    }
}

func (cas *CryptoApisService) GetNodeInfo() (CryptoApisGetNodeInfoRes, error) {
    getNodeInfoRes := CryptoApisGetNodeInfoRes{}

    restRes, err := resty.R().
      SetHeader("Content-Type", "application/json").
      SetHeader("X-API-Key", cas.ApiKey).
      Get("https://api.cryptoapis.io/v1/bc/btc/"+cas.Network+"/info")

    err = json.Unmarshal(restRes, &getNodeInfoRes)

    return getNodeInfoRes, err
}


