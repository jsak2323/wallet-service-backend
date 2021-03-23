package handlers

import (
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"

    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
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

func GetMaintenanceList(systemConfigRepo sc.SystemConfigRepository) (map[string]bool, error) {
    maintenanceList := map[string]bool{}
    maintenanceListObj, err := systemConfigRepo.GetByName(sc.MAINTENANCE_LIST)
    if err != nil { return maintenanceList, err }

    if maintenanceListObj.Value == "" {
        return maintenanceList, nil
    }

    maintenanceListSlice := strings.Split(maintenanceListObj.Value, ",")
    for _, symbol := range maintenanceListSlice {
        maintenanceList[symbol] = true
    }
    return maintenanceList, nil
}


