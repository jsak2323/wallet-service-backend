package auth

import(
    "time"
    "strings"
    "strconv"
    "crypto/md5"
    "crypto/sha256"
    "encoding/hex"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
)

func GenerateHashkey(rpcConfig rc.RpcConfig) (digestSha256String string, nonce string) {
    mt    := Microtime()
    nonce = strings.ReplaceAll(strconv.FormatFloat(mt, 'f', 9, 64), ".", "")

    unixTime := time.Now().Unix()
    this15m  := unixTime / 60

    pass    := rpcConfig.Password
    hashkey := rpcConfig.Hashkey

    digest := pass + strconv.FormatInt(this15m, 10) + hashkey + nonce
    digestMd5       := md5.Sum([]byte(digest))
    digestMd5String := hex.EncodeToString(digestMd5[:])
    digestSha256       := sha256.Sum256([]byte(digestMd5String))
    digestSha256String = hex.EncodeToString(digestSha256[:])

    return
}