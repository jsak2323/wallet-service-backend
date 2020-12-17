package cron

import(
    "fmt"
    "time"
    "strconv"
    "net/http"

    h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
    logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
    "github.com/btcid/wallet-services-backend-go/cmd/config"
    "github.com/btcid/wallet-services-backend-go/pkg/modules"
    "github.com/btcid/wallet-services-backend-go/pkg/lib/util"
)

type HealthCheckService struct {
    healthCheckRepo hc.HealthCheckRepository
    moduleServices *modules.ModuleServiceMap
}

func NewHealthCheckService(healthCheckRepo hc.HealthCheckRepository, moduleServices *modules.ModuleServiceMap) *HealthCheckService{
    return &HealthCheckService{
        healthCheckRepo,
        moduleServices,
    }
}

func (hcs *HealthCheckService) HealthCheckHandler(w http.ResponseWriter, req *http.Request) {
    isPing := false

    gbcRES := make(h.GetBlockCountHandlerResponseMap)
    getBlockCountService := h.NewGetBlockCountService(hcs.moduleServices)

    time.Sleep(time.Second*5)

    // after 9 or more minutes, save health check to db. otherwise, only ping
    lastHealthCheck, _ := hcs.healthCheckRepo.GetByRpcConfigId(1)
    minuteDiff, err := util.GetMinuteDiffFromNow(lastHealthCheck.LastUpdated)
    if err == nil && lastHealthCheck.Id == 1 && minuteDiff < 9 {
        isPing = true
    }

    logger.InfoLog(" - HealthCheckHandler Getting node blockcounts ..." , req)
    getBlockCountService.InvokeGetBlockCount(&gbcRES, "")
    logger.InfoLog(" - HealthCheckHandler Getting node blockcounts done. Fetched "+strconv.Itoa(len(gbcRES))+" results." , req)

    for resSymbol, resRpcConfigs := range gbcRES { 
        for _, resRpcConfig := range resRpcConfigs { 
            nodeBlockCount, _ := strconv.Atoi(resRpcConfig.Blocks)

            isBlockCountHealthy, blockDiff := false, 0

            if resRpcConfig.RpcConfig.IsHealthCheckEnabled {

                if isPing { // if ping, only check if blockount is 0
                    if nodeBlockCount <= 0 {
                        hcs.sendNotificationEmails(resRpcConfig)
                    }
                    fmt.Println(" -- Healthcheck ping "+resSymbol+" Blocks: "+resRpcConfig.Blocks)
                    continue
                }

                _isBlockCountHealthy, _blockDiff, err := (*hcs.moduleServices)[resSymbol].IsBlockCountHealthy(nodeBlockCount, resRpcConfig.RpcConfig.RpcConfigId)
                if err != nil { logger.ErrorLog(" - HealthCheckHandler hcs.ModuleServices[resSymbol].IsBlockCountHealthy(resRpcConfig.Blocks) err: "+err.Error()) }

                isBlockCountHealthy = _isBlockCountHealthy
                blockDiff           = _blockDiff

                if !isBlockCountHealthy && config.FirstHealthCheck { // if not healthy, send notification emails
                    hcs.sendNotificationEmails(resRpcConfig)
                }

                config.FirstHealthCheck = true
            }

            if !isPing {
                hcs.saveHealthCheck(resRpcConfig.RpcConfig.RpcConfigId, nodeBlockCount, blockDiff, isBlockCountHealthy)
            }
        }
    }
}

func (hcs *HealthCheckService) saveHealthCheck(rpcConfigId int, blockCount int, blockDiff int, isBlockCountHealthy bool) error {
    existingHealthCheck, err := hcs.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
    if err != nil { return err }

    if existingHealthCheck.Id == 0 { // does not exist, create a new one
        newHealthCheck := hc.HealthCheck{
            RpcConfigId         : rpcConfigId,
            BlockCount          : blockCount,
            BlockDiff           : blockDiff,
            IsHealthy           : isBlockCountHealthy,
        }
        err := hcs.healthCheckRepo.Create(&newHealthCheck)
        if err != nil {
            logger.ErrorLog(" - saveHealthCheck err: "+err.Error())
        } else {
            logger.Log(" - saveHealthCheck Create rpcConfigId: "+strconv.Itoa(newHealthCheck.Id)+" Success, HealthCheck: "+fmt.Sprintf("%+v", newHealthCheck))
        }

    } else { // already exists, update
        newHealthCheck := hc.HealthCheck{
            Id                  : existingHealthCheck.Id,
            RpcConfigId         : rpcConfigId,
            BlockCount          : blockCount,
            BlockDiff           : blockDiff,
            IsHealthy           : isBlockCountHealthy,
        }
        err := hcs.healthCheckRepo.Update(&newHealthCheck)
        if err != nil {
            logger.ErrorLog(" - saveHealthCheck err: "+err.Error())
        } else {
            logger.Log(" - saveHealthCheck Update rpcConfigId: "+strconv.Itoa(rpcConfigId)+" Success, HealthCheck: "+fmt.Sprintf("%+v", newHealthCheck))
        }
    }

    return nil
}

func (hcs *HealthCheckService) sendNotificationEmails(res h.GetBlockCountRes) {
    logger.Log(" - HealthCheckHandler -- Sending notification email ...")

    subject := "Health Check Failed for "+res.RpcConfig.Symbol+" VM ("+res.RpcConfig.Host+")"
    message := "Health check has failed with following detail: "+
    "\n Symbol: "+res.RpcConfig.Symbol+
    "\n Host: "+res.RpcConfig.Host+
    "\n Name: "+res.RpcConfig.Name+
    "\n Type: "+res.RpcConfig.Type+
    "\n Node Version: "+res.RpcConfig.NodeVersion+
    "\n BlockCount: "+res.Blocks

    recipients := config.CONF.NotificationEmails

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog(err.Error()) }
    logger.Log(" - HealthCheckHandler -- Is notification email sent: "+strconv.FormatBool(isEmailSent))
}


