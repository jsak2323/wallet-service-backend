package fireblocks

import (
	"errors"
	"encoding/json"
	
	"gopkg.in/resty.v0"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
)

const getAssetEndpoint = "getVaultAccountAsset"

func GetVaultAccountAsset(req GetVaultAccountAssetReq) (RES GetVaultAccountAssetRes, err error) {
	res, err := resty.R().
		SetHeaders(
			map[string]string{
				"Content-type": "application/json",
				"Authorization": "Basic " + auth(),
			},
		).
		Get(config.CONF.FireblocksHost+"/"+getAssetEndpoint+"/"+req.VaultAccountId+"/"+req.AssetId)

    if err != nil {
        return GetVaultAccountAssetRes{}, err
    }

	if err = json.Unmarshal(res.Body, &RES); err != nil {
		return GetVaultAccountAssetRes{}, err
	}

	if RES.Error != "" {
		return GetVaultAccountAssetRes{}, errors.New(RES.Error)
	}

	return RES, nil
}