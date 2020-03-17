package util

import(
    "bytes"
    "time"
    "net/http"

    "github.com/divan/gorilla-xmlrpc/xml"

    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
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

func NewXmlRpc(host string, port string, path string) *XmlRpc{
    return &XmlRpc{
        host, port, path,
    }
}

func (xr *XmlRpc) XmlRpcCall(method string, args *RpcReq, reply interface{}) error {
    buf, err := xml.EncodeClientRequest(method, args)
    if err != nil { return err }

    url := "http://"+xr.Host+":"+xr.Port+xr.Path
    httpClient := &http.Client{
        Timeout: 5 * time.Second,
    }
    resp, err := httpClient.Post(url, "text/xml", bytes.NewBuffer(buf))
    if err != nil { return err }
    
    defer resp.Body.Close()

    err = xml.DecodeClientResponse(resp.Body, reply)
    if err != nil { return err }
    
    return nil
}

func GenerateRpcReq(rpcConfig rc.RpcConfig, arg1 string, arg2 string, arg3 string) RpcReq {
    hashkey, nonce := generateHashkey(rpcConfig)
    
    return RpcReq{
        RpcUser : rpcConfig.User,
        Hash    : hashkey,
        Arg1    : arg1,
        Arg2    : arg2,
        Arg3    : arg3,
        Nonce   : nonce,
    }
}