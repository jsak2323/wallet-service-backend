package xmlrpc

import (
    // "fmt"
    "os"
    "io"
    "time"
    "net/http"

    "github.com/cactus/gostrftime"

    // logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
)

func (gs *GeneralService) GetLog(rpcConfig rc.RpcConfig, date string) (string, error) {
    // generate filepath
    filepath := "download/"+gs.Symbol+"-"+rpcConfig.Type+"-"+date
    now := time.Now()
    if gostrftime.Format("%m-%d-%Y", now) == date { // if today, include time
        timestamp := gostrftime.Format("%H:%M:%S", now)
        filepath += "-"+timestamp
    }
    filepath += ".log"

    // fetch log file data
    res, err := http.Get("http://"+rpcConfig.Host+":"+rpcConfig.Port+"/log/"+date)
    if err != nil { return "", err }
    defer res.Body.Close()

    // create file
    out, err := os.Create(filepath)
    if err != nil { return "", err }
    defer out.Close()

    // write body to file
    _, err = io.Copy(out, res.Body)
    if err != nil { return "", err }

    return filepath, nil
}


