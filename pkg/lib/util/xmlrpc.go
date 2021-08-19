package util

import (
    "fmt"
    "bytes"
    "time"
    "strings"
    "strconv"
    "net/http"
    "crypto/md5"
    "crypto/sha256"
    "encoding/hex"

    "github.com/divan/gorilla-xmlrpc/xml"

    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
)

type RpcReq struct {
    RpcUser string
    Hash    string
    Arg1    string
    Arg2    string
    Arg3    string
    Nonce   string
}

type XmlRpc struct {
    Host string
    Port string
    Path string
}
func NewXmlRpcClient(host string, port string, path string) *XmlRpc{
    return &XmlRpc{
        host, port, path,
    }
}
func (xr *XmlRpc) XmlRpcCall(method string, args *RpcReq, reply interface{}) error {
    buf, err := xml.EncodeClientRequest(method, args)
    if err != nil { 
        fmt.Println(" - xml.EncodeClientRequest(method, args) err: "+err.Error())
        return err 
    }

    url := "http://"+xr.Host+":"+xr.Port+xr.Path
    httpClient := &http.Client{
        Timeout: 120 * time.Second,
    }
    res, err := httpClient.Post(url, "text/xml", bytes.NewBuffer(buf))
    if err != nil { 
        fmt.Println(" - httpClient.Post xml err: "+err.Error())
        return err 
    }

    defer res.Body.Close()
    
    err = xml.DecodeClientResponse(res.Body, reply)
    if err != nil { 
        fmt.Println(" - xml.DecodeClientResponse err: "+err.Error())
        return err 
    }
    
    return nil
}

func GenerateRpcReq(rpcConfig rc.RpcConfig, arg1 string, arg2 string, arg3 string) RpcReq {
    hashkey, nonce := GenerateHashkey(rpcConfig.Password, rpcConfig.Hashkey)
    
    return RpcReq{
        RpcUser : rpcConfig.User,
        Hash    : hashkey,
        Arg1    : arg1,
        Arg2    : arg2,
        Arg3    : arg3,
        Nonce   : nonce,
    }
}

func GenerateHashkey(pass, hashkey string) (digestSha256String string, nonce string) {
    mt    := Microtime()
    nonce = strings.ReplaceAll(strconv.FormatFloat(mt, 'f', 9, 64), ".", "")

    unixTime := time.Now().Unix()
    this15m  := unixTime / 60

    // todo: implement encryption
    // pass    := rpcConfig.Password
    // hashkey := rpcConfig.Hashkey

    digest := pass + strconv.FormatInt(this15m, 10) + hashkey + nonce
    digestMd5       := md5.Sum([]byte(digest))
    digestMd5String := hex.EncodeToString(digestMd5[:])
    digestSha256       := sha256.Sum256([]byte(digestMd5String))
    digestSha256String = hex.EncodeToString(digestSha256[:])

    return
}


