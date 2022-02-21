package walletrpc

import (
	"context"
	"strings"

	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

func (s *walletRpcService) ListMaintenance(ctx context.Context, action, value string) (err error) {

	// get maintenance list
	maintenanceList, err := s.GetMaintenanceList(ctx)
	if err != nil {
		err = errs.AssignErr(errs.AddTrace(err), errs.FailedGetMaintenanceList)
		return err
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

	if (action == "remove" && maintenanceList[value]) || len(symbolArray) > 0 {
		updateValue := strings.Join(symbolArray, ",")

		updateErr := s.systemConfigRepo.Update(sc.SystemConfig{
			Name:  sc.MAINTENANCE_LIST,
			Value: updateValue,
		})
		if updateErr != nil {
			err = errs.AssignErr(errs.AddTrace(err), errs.FailedUpdateSystemConfig)
			return err
		}
	}

	return err
}
