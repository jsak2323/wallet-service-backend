package handlers

import(
    "fmt"
    "bytes"
    "net/http"

    "github.com/btcid/wallet-services-backend/pkg/lib/util"

    "github.com/gorilla/mux"
)

type GetBlockCountRes struct {
    Blocks string
}

func GetBlockCountHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    w.WriteHeader(http.StatusOK)
    // fmt.Fprintf(w, "Symbol: %v\n", vars["symbol"])
    fmt.Println("Symbol: "+vars["symbol"])

    var res GetBlockCountRes

    rpcReq := RpcReq{
        RpcUser     : "testuser",
        Hash        : "hash",
        Arg1        : "arg1",
        Arg2        : "arg2",
        Arg3        : "arg3",
        Nonce       : "nonce",
    }

    xmlrpc := util.NewXmlRpc("35.187.234.25", "3000", "/rpc")
    xmlrpc.XmlRpcCall("EthRpc.GetBlockCount", &rpcReq, &res)

    fmt.Println("res: ") 
    fmt.Printf("%+v", res)

}