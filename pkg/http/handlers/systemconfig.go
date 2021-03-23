package handlers

import (
    "strings"
    "net/http"
    "encoding/json"

    "github.com/gorilla/mux"

    sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type SystemConfigService struct {
    systemConfigRepo sc.SystemConfigRepository
}

func NewSystemConfigService(systemConfigRepo sc.SystemConfigRepository) *SystemConfigService {
    return &SystemConfigService{
        systemConfigRepo,
    }
}

func (scs *SystemConfigService) MaintenanceListHandler(w http.ResponseWriter, req *http.Request) { 
    vars := mux.Vars(req)
    action := vars["action"]
    value  := strings.ToUpper(vars["value"])

    // define response object
    RES := StandardRes{}

    // define response handler
    handleResponse := func() {
        resStatus := http.StatusOK
        if RES.Error != "" {
            resStatus = http.StatusInternalServerError
        }
        w.WriteHeader(resStatus)
        json.NewEncoder(w).Encode(RES)
    }
    defer handleResponse()

    // get maintenance list
    maintenanceList, err := GetMaintenanceList(scs.systemConfigRepo)
    if err != nil { 
        logger.ErrorLog(" - MaintenanceListHandler GetMaintenanceList err: "+err.Error()) 
        RES.Error = err.Error()
        return
    }

    symbolArray := []string{}
    if action == "add" {
        if !maintenanceList[value] {
            for SYMBOL, _ := range maintenanceList {
                symbolArray = append(symbolArray, SYMBOL)
            }
            symbolArray = append(symbolArray, value)
        }

    } else if action == "remove" {
        if maintenanceList[value] {
            for SYMBOL, _ := range maintenanceList {
                if SYMBOL != value {
                    symbolArray = append(symbolArray, SYMBOL)
                }
            }
        }
    }

    updateValue := strings.Join(symbolArray, ",")
    err := scs.systemConfigRepo.Update(sc.SystemConfig{
        Name  : sc.MAINTENANCE_LIST,
        Value : updateValue,
    })
    if err != nil {
        logger.ErrorLog(" - MaintenanceListHandler scs.systemConfigRepo.Update err: "+err.Error()) 
        RES.Error = err.Error()
        return
    }

    // handle success response
    RES.Success = true
    resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - MaintenanceListHandler Success. Res: "+string(resJson), req)
}


