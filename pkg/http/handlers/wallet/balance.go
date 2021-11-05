package wallet

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	modulesm "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type GetBalanceHandlerResponseMap map[int]GetBalanceRes

func (s *WalletService) GetBalanceHandler(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
    currencyId, err := strconv.Atoi(vars["currency_id"])
	if err != nil {

	}
    isGetAll := currencyId == 0

	RES := make(GetBalanceHandlerResponseMap)

    if isGetAll {
        logger.InfoLog(" - wallet.GetBalanceHandler For all currencies, Requesting ...", req) 
    } else {
        logger.InfoLog(" - wallet.GetBalanceHandler For currency_id: "+strconv.Itoa(currencyId)+", Requesting ...", req) 
    }

	s.InvokeGetBalance(&RES, currencyId)
	
	resJson, _ := json.Marshal(RES)
    logger.InfoLog(" - wallet.GetBalanceHandler Success. CurrencyId: "+strconv.Itoa(currencyId)+", Res: "+string(resJson), req)
	w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(RES)
}

func (s *WalletService) InvokeGetBalance(RES *GetBalanceHandlerResponseMap, currencyId int) {
	if currencyId == 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(config.CURRRPC))

		for _, curr := range config.CURRRPC {
			go func(currencyConfiguration config.CurrencyRpcConfig) {
				defer wg.Done()
				
				(*RES)[currencyId] = s.GetBalance(currencyConfiguration)
			}(curr)
		}

		wg.Wait()
    } else {
		(*RES)[currencyId] = s.GetBalance(config.CURRRPC[currencyId])
	}
}

func (s *WalletService) GetBalance(currConfig config.CurrencyRpcConfig) GetBalanceRes {
	var wg sync.WaitGroup
	var res GetBalanceRes =  GetBalanceRes{CurrencyConfig: currConfig.Config}
	
	wg.Add(5)
	go func() { defer wg.Done(); s.SetColdBalanceDetails(&res) }()
	go func() { defer wg.Done(); s.SetHotBalanceDetails(currConfig.RpcConfigs, &res) }()
	go func() { defer wg.Done(); s.SetUserBalanceDetails(&res) }()
	go func() { defer wg.Done(); s.SetPendingWithdraw(&res) }()
	go func() { defer wg.Done(); s.SetHotLimits(&res) }()
	wg.Wait()

	s.SetPercent(&res)
	
	return res
}

func (s *WalletService) SetColdBalanceDetails(res *GetBalanceRes) {
	var symbol string = res.CurrencyConfig.Symbol
	var cbs    []cb.ColdBalance = s.coldWalletService.GetBalance(res.CurrencyConfig.Id)
	var err    error
	
	for _, cb := range cbs {
		var coldBalanceDetail = BalanceDetail{ Id: cb.Id, Name: cb.Name, Type: cb.Type }

		coldBalanceDetail.Coin = cb.Balance
		coldBalanceDetail.Address = cb.Address
		coldBalanceDetail.FireblocksName = cb.FireblocksName
		
		if coldBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(coldBalanceDetail.Coin, symbol); err != nil {
			logger.ErrorLog(" - SetColdBalanceDetails ConvertCoinToIdr("+cb.Type+", "+cb.Balance+") err: "+err.Error())
		}

		if res.TotalColdCoin, err = util.AddCoin(res.TotalColdCoin, coldBalanceDetail.Coin); err != nil {
			logger.ErrorLog(" - SetColdBalanceDetails AddCoin("+cb.Type+", "+cb.Balance+") err: "+err.Error())
		}

		if res.TotalColdIdr, err = util.AddIdr(res.TotalColdIdr, coldBalanceDetail.Idr); err != nil {
			logger.ErrorLog(" - SetColdBalanceDetails AddIdr("+cb.Type+", "+cb.Balance+") err: "+err.Error())
		}

		res.ColdBalances = append(res.ColdBalances, coldBalanceDetail)
	}
}

func (s *WalletService) SetHotBalanceDetails(rpcConfigs []rc.RpcConfig, res *GetBalanceRes) {
	var symbol string = res.CurrencyConfig.Symbol
	
	for _, rpcConfig := range rpcConfigs {
		var hotBalanceDetail BalanceDetail = BalanceDetail{ Name: rpcConfig.Name, Type: rpcConfig.Type }
		var rpcRes 			 *modulesm.GetBalanceRpcRes

		module, err := s.moduleServices.GetModule(res.CurrencyConfig.Id);
		if err != nil {
			logger.ErrorLog(" - s.moduleServices err: "+err.Error())
			continue
		}

		if rpcRes, err = module.GetBalance(rpcConfig); err != nil {
			logger.ErrorLog(" - SetHotBalanceDetails node.GetBalance("+symbol+", "+rpcConfig.Name+") err: "+err.Error())
			continue
		}

		hotBalanceDetail.Coin = rpcRes.Balance

		if hotBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(hotBalanceDetail.Coin, symbol); err != nil {
			logger.ErrorLog(" - SetHotBalanceDetails node.ConvertCoinToIdr("+symbol+", "+rpcConfig.Name+") err: "+err.Error())
		}

		if res.TotalHotCoin, err = util.AddCoin(res.TotalHotCoin, hotBalanceDetail.Coin); err != nil {
			logger.ErrorLog(" - SetHotBalanceDetails AddCoin("+symbol+", "+rpcConfig.Name+") err: "+err.Error())
		}

		if res.TotalHotIdr, err = util.AddIdr(res.TotalHotIdr, hotBalanceDetail.Idr); err != nil {
			logger.ErrorLog(" - SetHotBalanceDetails AddIdr("+symbol+", "+rpcConfig.Name+") err: "+err.Error())
		}

		res.HotBalances = append(res.HotBalances, hotBalanceDetail)
	}
}

func (s *WalletService) SetUserBalanceDetails(res *GetBalanceRes) {
	var tcb    ub.TotalCoinBalance
	var err    error
	var symbol string = res.CurrencyConfig.Symbol
	var frozenBalanceDetail BalanceDetail = BalanceDetail{ Name: "Frozen" }
	var liquidBalanceDetail BalanceDetail = BalanceDetail{ Name: "Liquid" }
	
	if tcb, err = s.userBalanceRepo.GetTotalCoinBalance(symbol); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails ub.GetTotalCoinBalance("+symbol+") err: "+err.Error())
	}

	if liquidBalanceDetail.Coin, err = util.RawToCoin(tcb.Total, 8); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails RawToCoin("+symbol+") err: "+err.Error())
	} else if liquidBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(liquidBalanceDetail.Coin, symbol); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails liquid.ConvertCoinToIdr("+symbol+") err: "+err.Error())
	}
	
	res.UserBalances = append(res.UserBalances, liquidBalanceDetail)
	
	if frozenBalanceDetail.Coin, err = util.RawToCoin(tcb.TotalFrozen, 8); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails RawToCoin("+symbol+") err: "+err.Error())
	} else if frozenBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(frozenBalanceDetail.Coin, symbol); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails frozen.ConvertCoinToIdr("+symbol+") err: "+err.Error())
	}

	res.UserBalances = append(res.UserBalances, frozenBalanceDetail)

	if res.TotalUserCoin, err = util.AddCoin(liquidBalanceDetail.Coin, frozenBalanceDetail.Coin); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails AddCoin("+symbol+") err: "+err.Error())
	}

	if res.TotalUserIdr, err = util.AddIdr(liquidBalanceDetail.Idr, frozenBalanceDetail.Idr); err != nil {
		logger.ErrorLog(" - SetUserBalanceDetails AddIdr("+symbol+") err: "+err.Error())
	}
}

func (s *WalletService) SetPendingWithdraw(res *GetBalanceRes) {
	var err error
	var symbol string = res.CurrencyConfig.Symbol
	var pendingWDRaw string

	if pendingWDRaw, err = s.withdrawRepo.GetPendingWithdraw(symbol); err != nil {
		logger.ErrorLog(" - SetPendingWithdraw GetPendingWithdraw("+symbol+") err: "+err.Error())
		return
	}

	if res.PendingWDCoin, err = util.RawToCoin(pendingWDRaw, 8); err != nil {
		logger.ErrorLog(" - SetPendingWithdraw RawToCoin("+symbol+", "+pendingWDRaw+") err: "+err.Error())
	} else if res.PendingWDIdr, err = s.marketService.ConvertCoinToIdr(res.PendingWDCoin, symbol); err != nil {
		logger.ErrorLog(" - SetPendingWithdraw ConvertCoinToIdr("+symbol+") err: "+err.Error())
	}
}

func (s *WalletService) SetPercent(res *GetBalanceRes) {
	var err 	error
	var hotCold string
	
	if res.HotPercent, err = util.PercentBig(res.TotalHotCoin, res.TotalUserCoin); err != nil {
		logger.ErrorLog(" - SetPercent PercentBig("+res.CurrencyConfig.Symbol+") err: "+err.Error())
	}

	if hotCold, err = util.AddCoin(res.TotalColdCoin, res.TotalHotCoin); err != nil {
		logger.ErrorLog(" - SetPercent AddCoin("+res.CurrencyConfig.Symbol+") err: "+err.Error())
	}

	if res.HotColdPercent, err = util.PercentBig(hotCold, res.TotalUserCoin); err != nil {
		logger.ErrorLog(" - SetPercent PercentBig("+res.CurrencyConfig.Symbol+") err: "+err.Error())
	}
}

func (s *WalletService) SetHotLimits(res *GetBalanceRes) {
	var err error
	
	if res.HotLimits, err = s.hotLimitRepo.GetBySymbol(res.CurrencyConfig.Symbol); err != nil {
		logger.ErrorLog(" - SetHotLimits hotLimitRepo.GetBySymbol("+res.CurrencyConfig.Symbol+") err: "+err.Error())	
	}
}
