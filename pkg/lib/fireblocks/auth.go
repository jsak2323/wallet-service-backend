package fireblocks

import (
	"encoding/base64"
	
	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func auth() string {
	hash, nonce := util.GenerateHashkey(config.CONF.FireblocksServerpass, config.CONF.FireblocksServerhashkey)
	
	data := []byte(config.CONF.FireblocksServeruser+":"+hash+":"+nonce)

	return base64.StdEncoding.EncodeToString(data)
}