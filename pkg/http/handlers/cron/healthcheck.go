package cron

import(
    "fmt"
    "strconv"
    "net/http"

    h "github.com/btcid/wallet-services-backend/pkg/http/handlers"
    logger "github.com/btcid/wallet-services-backend/pkg/logging"
    "github.com/btcid/wallet-services-backend/cmd/config"
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type HealthCheckService struct{
    healthCheckRepo hc.HealthCheckRepository
}

func NewHealthCheckService(healthCheckRepo hc.HealthCheckRepository) *HealthCheckService{
    return &HealthCheckService{
        healthCheckRepo,
    }
}

func (hcs *HealthCheckService) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {

    getBlockCountService := h.NewGetBlockCountService()

    // get node blockcounts 
    logger.InfoLog("HealthCheckHandler Getting node blockcounts ..." , r)
    gbcRES := make(h.GetBlockCountHandlerResponseMap)
    getBlockCountService.InvokeGetBlockCount(&gbcRES, "", false)
    logger.InfoLog("HealthCheckHandler Getting node blockcounts done. Fetched "+strconv.Itoa(len(gbcRES))+" results." , r)

    // get comparation blockcounts
    logger.InfoLog("HealthCheckHandler Getting comparation blockcounts ..." , r)
    cbcRES := make(h.GetBlockCountHandlerResponseMap)
    getBlockCountService.InvokeGetBlockCount(&cbcRES, "", true)
    logger.InfoLog("HealthCheckHandler Getting comparation blockcounts done. Fetched "+strconv.Itoa(len(cbcRES))+" results." , r)

    // compare the results
    for resSymbol, resRpcConfigs := range gbcRES {
        for _, resRpcConfig := range resRpcConfigs {        
            nodeBlocks          := resRpcConfig.Blocks
            comparationBlocks   := hcs.findComparationResultByRpcConfigId(resRpcConfig.RpcConfigId, &cbcRES).Blocks

            nodeBlocksInt, _        := strconv.Atoi(nodeBlocks)
            comparationBlocksInt, _ := strconv.Atoi(comparationBlocks)

            err := hcs.saveHealthCheck(resRpcConfig.RpcConfigId, nodeBlocksInt, comparationBlocksInt)
            if err != nil { logger.ErrorLog("HealthCheckHandler hcs.saveHealthCheck(resRpcConfig, nodeBlocks, comparationBlocks), err: "+err.Error()) }


            blocksDiff := nodeBlocksInt - comparationBlocksInt
            if blocksDiff >= config.CURR[resSymbol].Config.HealthyBlockDiff {
                fmt.Println("Block Difference detected.")
            }



        }
    }
}

func (hcs *HealthCheckService) findComparationResultByRpcConfigId(rpcConfigId int, RES *h.GetBlockCountHandlerResponseMap) h.GetBlockCountRes {
    emptyRes := h.GetBlockCountRes{}
    for _, resRpcConfigs := range (*RES) {
        for _, resRpcConfig := range resRpcConfigs {
            if resRpcConfig.RpcConfigId == rpcConfigId { return resRpcConfig }
        }
    }
    return emptyRes
}

func (hcs *HealthCheckService) saveHealthCheck(rpcConfigId int, blockCount int, confirmBlockCount int) error {
    
    existingHealthCheck, err := hcs.healthCheckRepo.GetByRpcConfigId(rpcConfigId)
    if err != nil { return err }

    if existingHealthCheck.Id == 0 { // does not exist, create a new one
        newHealthCheck := hc.HealthCheck{
            RpcConfigId         : rpcConfigId,
            BlockCount          : blockCount,
            ConfirmBlockCount   : confirmBlockCount,
        }
        err := hcs.healthCheckRepo.Create(&newHealthCheck)
        if err != nil {
            logger.ErrorLog("saveHealthCheck err: "+err.Error())
        } else {
            logger.Log("saveHealthCheck Success, HealthCheck: "+fmt.Sprintf("%+v", newHealthCheck))
        }

    } else { // already exists, update
        fmt.Println("update hit")
    }

    return nil
}