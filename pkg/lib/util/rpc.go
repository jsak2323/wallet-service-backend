package util

import(
    "bytes"
    "net/http"

    "github.com/divan/gorilla-xmlrpc/xml"
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
    resp, err := http.Post(url, "text/xml", bytes.NewBuffer(buf))
    if err != nil { return err }
    
    defer resp.Body.Close()

    err = xml.DecodeClientResponse(resp.Body, reply)
    if err != nil { return err }
    
    return nil
}