package handlers

import(
    "fmt"
    "net/http"
    "encoding/json"

    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    ethservice "github.com/btcid/wallet-services-backend/pkg/modules/eth"

    "github.com/gorilla/mux"
)

func GetBlockCountHandler(w http.ResponseWriter, r *http.Request) { 
    vars := mux.Vars(r)
    symbol := vars["symbol"]    

    var handleSuccess = func(res *ethservice.GetBlockCountRes) {
        logger.InfoLog("GetBlockCountHandler Success. Symbol: "+symbol+", Blocks: "+res.Blocks, r)
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(res)
    }

    var handleError = func(err error, funcName string) {
        errMsg := "GetBlockCountHandler "+funcName+" Error: "+err.Error()
        logger.ErrorLog(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
    }

    logger.InfoLog("GetBlockCountHandler Symbol: "+symbol+", Requesting ...", r) 

    switch symbol { 
        case "eth" :
            res, err := ethservice.GetBlockCount()
            if err != nil { 
                handleError(err, "ethservice.GetBlockCount()") 
                return
            }
            handleSuccess(res) 

        default :
            fmt.Println("default")
    }

}
