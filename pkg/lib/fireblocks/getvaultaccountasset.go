package fireblocks

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

const getAssetEndpoint = "getVaultAccountAsset"

func GetVaultAccountAsset(req GetVaultAccountAssetReq) (RES GetVaultAccountAssetRes, err error) {
	httpClient := &http.Client{
        Timeout: 120 * time.Second,
    }

	res, err := httpClient.Get(config.CONF.FireblocksHost+"/"+getAssetEndpoint+"/"+req.VaultAccountId+"/"+req.AssetId)
    if err != nil {
        return GetVaultAccountAssetRes{}, err
    }
    defer res.Body.Close()

	if err = json.NewDecoder(res.Body).Decode(&RES); err != nil {
		return GetVaultAccountAssetRes{}, err
	}

	return RES, nil
}