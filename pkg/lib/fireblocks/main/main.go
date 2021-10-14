package main

import (
	"encoding/base64"
	"fmt"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

func main() {
	hash, nonce := util.GenerateHashkey("sr8|1o6B~li3W_.", "ds9FntYP55")
	
	data := []byte("fireblocks"+":"+hash+":"+nonce)

	fmt.Println("______", base64.StdEncoding.EncodeToString(data))
	fmt.Println("______", config.CONF.FireblocksServerhashkey)
	fmt.Println("______", hash)

}