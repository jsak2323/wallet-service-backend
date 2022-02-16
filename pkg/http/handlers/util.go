package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func DecodeAndLogPostRequest(req *http.Request, output interface{}) error {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errs.AddTrace(err)
	}

	logger.InfoLog("POST Request Body : "+string(reqBody), req)

	err = json.Unmarshal(reqBody, output)
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}

func GetMaintenanceList(ctx context.Context, systemConfigRepo sc.Repository) (map[string]bool, error) {
	maintenanceList := map[string]bool{}
	maintenanceListObj, err := systemConfigRepo.GetByName(ctx, sc.MAINTENANCE_LIST)
	if err != nil {
		return maintenanceList, errs.AddTrace(err)
	}

	if maintenanceListObj.Value == "" {
		return maintenanceList, nil
	}

	maintenanceListSlice := strings.Split(maintenanceListObj.Value, ",")
	for _, symbol := range maintenanceListSlice {
		maintenanceList[symbol] = true
	}
	return maintenanceList, nil
}
