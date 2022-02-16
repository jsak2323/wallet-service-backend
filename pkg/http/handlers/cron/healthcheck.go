package cron

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	sc "github.com/btcid/wallet-services-backend-go/pkg/domain/systemconfig"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type HealthCheckService struct {
	moduleServices   *modules.ModuleServiceMap
	healthCheckRepo  hc.Repository
	systemConfigRepo sc.Repository
}

func NewHealthCheckService(
	moduleServices *modules.ModuleServiceMap,
	healthCheckRepo hc.Repository,
	systemConfigRepo sc.Repository,
) *HealthCheckService {
	return &HealthCheckService{
		moduleServices,
		healthCheckRepo,
		systemConfigRepo,
	}
}

func (hcs *HealthCheckService) HealthCheckHandler() {

	var (
		startTime = time.Now()
		isPing    = false
		gbcRES    = make(h.GetBlockCountHandlerResponseMap)

		errField *errs.Error = nil
		ctx                  = context.Background()
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	getBlockCountService := h.NewGetBlockCountService(hcs.moduleServices, hcs.systemConfigRepo)

	// after 11 or more minutes, save health check to db. otherwise, only ping
	lastUpdatedTime, _ := hcs.healthCheckRepo.GetLastUpdatedTime()
	minuteDiff, err := util.GetMinuteDiffFromNow(lastUpdatedTime)
	if err == nil && minuteDiff < 11 {
		isPing = true
	}

	// get maintenance list
	maintenanceList, err := h.GetMaintenanceList(ctx, hcs.systemConfigRepo)
	if err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.ErrorHealthCheckHandler)
	}

	logger.Log(" - HealthCheckHandler Getting node blockcounts ...", ctx)
	getBlockCountService.InvokeGetBlockCount(ctx, &gbcRES, "", "")
	logger.Log(" - HealthCheckHandler Getting node blockcounts done. Fetched "+strconv.Itoa(len(gbcRES))+" results.", ctx)

	for resSymbol, mapTokenType := range gbcRES {
		for resTokenType, resRpcConfigs := range mapTokenType {
			for _, resRpcConfig := range resRpcConfigs {
				nodeBlockCount, _ := strconv.Atoi(resRpcConfig.Blocks)

				isBlockCountHealthy, blockDiff := false, 0

				if !resRpcConfig.RpcConfig.IsHealthCheckEnabled {
					continue
				}
				if maintenanceList[resSymbol] {
					continue
				}

				if isPing { // if ping, only check if blockount is 0
					if nodeBlockCount <= 0 {
						hcs.sendNotificationEmails(ctx, resRpcConfig, -1)
					}
					fmt.Println(" -- Healthcheck ping " + resSymbol + " Blocks: " + resRpcConfig.Blocks)
					continue
				}

				currencyConfig, err := config.GetCurrencyBySymbolTokenType(resSymbol, resTokenType)
				if err != nil {
					errField = errs.AssignErr(errs.AddTrace(err), errs.ErrorHealthCheckHandler)
					continue
				}

				module, err := hcs.moduleServices.GetModule(currencyConfig.Id)
				if err != nil {
					errField = errs.AssignErr(errs.AddTrace(err), errs.ErrorHealthCheckHandler)
					continue
				}

				_isBlockCountHealthy, _blockDiff, err := module.IsBlockCountHealthy(ctx, nodeBlockCount, resRpcConfig.RpcConfig.RpcConfigId)
				if err != nil {
					errField = errs.AssignErr(errs.AddTrace(err), errs.ErrorHealthCheckHandler)
				}

				isBlockCountHealthy = _isBlockCountHealthy
				blockDiff = _blockDiff

				if !isBlockCountHealthy && config.FirstHealthCheck { // if not healthy, send notification emails
					hcs.sendNotificationEmails(ctx, resRpcConfig, blockDiff)
				}

				config.FirstHealthCheck = true

				if !isPing {
					hcs.saveHealthCheck(ctx, resRpcConfig.RpcConfig.RpcConfigId, nodeBlockCount, blockDiff, isBlockCountHealthy)
				}
			}
		}
	}

	elapsedTime := time.Since(startTime)
	fmt.Println(" - HealthCheckHandler Time elapsed: " + fmt.Sprintf("%f", elapsedTime.Minutes()) + " minutes.")
}

func (hcs *HealthCheckService) saveHealthCheck(ctx context.Context, rpcConfigId int, blockCount int, blockDiff int, isBlockCountHealthy bool) error {
	existingHealthCheck, err := hcs.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
	if err != nil {
		return errs.AddTrace(err)
	}

	if existingHealthCheck.Id == 0 { // does not exist, create a new one
		newHealthCheck := hc.HealthCheck{
			RpcConfigId: rpcConfigId,
			BlockCount:  blockCount,
			BlockDiff:   blockDiff,
			IsHealthy:   isBlockCountHealthy,
		}
		err := hcs.healthCheckRepo.Create(&newHealthCheck)
		if err != nil {
			err = errs.AddTrace(err)
		} else {
			logger.Log(" - saveHealthCheck Create rpcConfigId: "+strconv.Itoa(newHealthCheck.Id)+" Success, HealthCheck: "+fmt.Sprintf("%+v", newHealthCheck), ctx)
		}

	} else { // already exists, update
		newHealthCheck := hc.HealthCheck{
			Id:          existingHealthCheck.Id,
			RpcConfigId: rpcConfigId,
			BlockCount:  blockCount,
			BlockDiff:   blockDiff,
			IsHealthy:   isBlockCountHealthy,
		}
		err := hcs.healthCheckRepo.Update(&newHealthCheck)
		if err != nil {
			err = errs.AddTrace(err)
		} else {
			logger.Log(" - saveHealthCheck Update rpcConfigId: "+strconv.Itoa(rpcConfigId)+" Success, HealthCheck: "+fmt.Sprintf("%+v", newHealthCheck), ctx)
		}
	}

	return nil
}

func (hcs *HealthCheckService) sendNotificationEmails(ctx context.Context, res h.GetBlockCountRes, blockDiff int) {
	logger.Log(" - HealthCheckHandler -- Sending notification email ...", ctx)

	blockCount := res.Blocks
	if blockCount == "" {
		blockCount = "0"
	}

	subject := "Health Check Failed for " + res.RpcConfig.Symbol + " VM (" + res.RpcConfig.Host + ")"
	message := "Health check has failed with following detail: " +
		"\n Symbol: " + res.RpcConfig.Symbol +
		"\n Host: " + res.RpcConfig.Host +
		"\n Name: " + res.RpcConfig.Name +
		"\n Type: " + res.RpcConfig.Type +
		"\n Node Version: " + res.RpcConfig.NodeVersion +
		"\n BlockCount: " + blockCount +
		"\n BlockDiff: " + strconv.Itoa(blockDiff)

	recipients := config.CONF.NotificationEmails

	isEmailSent, err := util.SendEmail(subject, message, recipients)
	if err != nil {
		logger.ErrorLog(err.Error(), context.Background())
	}
	logger.Log(" - HealthCheckHandler -- Is notification email sent: "+strconv.FormatBool(isEmailSent), ctx)
}
