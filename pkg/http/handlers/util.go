package handlers

import (
    "net/http"
    "io/ioutil"
    "encoding/json"

    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func DecodeAndLogPostRequest(req *http.Request, output interface{}) error {
    reqBody, err := ioutil.ReadAll(req.Body)
    if err != nil { return err }

    logger.InfoLog("POST Request Body : "+string(reqBody), req)

    err = json.Unmarshal(reqBody, output)
    if err != nil { return err }

    return nil
}


