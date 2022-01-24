package wallet

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	ub "github.com/btcid/wallet-services-backend-go/pkg/domain/userbalance"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	modulesm "github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

type GetBalanceHandlerResponseMap map[int]GetBalanceRes

func (s *WalletService) GetBalanceHandler(w http.ResponseWriter, req *http.Request) {
	var (
		vars = mux.Vars(req)

		ctx = req.Context()
	)
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

	s.InvokeGetBalance(ctx, &RES, currencyId)

	resJson, _ := json.Marshal(RES)
	logger.InfoLog(" - wallet.GetBalanceHandler Success. CurrencyId: "+strconv.Itoa(currencyId)+", Res: "+string(resJson), req)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RES)
}

func (s *WalletService) InvokeGetBalance(ctx context.Context, RES *GetBalanceHandlerResponseMap, currencyId int) {
	if currencyId == 0 {
		wg := sync.WaitGroup{}
		wg.Add(len(config.CURRRPC))

		for _, curr := range config.CURRRPC {
			go func(currencyConfiguration config.CurrencyRpcConfig) {
				defer wg.Done()

				(*RES)[currencyId] = s.GetBalance(ctx, currencyConfiguration)
			}(curr)
		}

		wg.Wait()
	} else {
		(*RES)[currencyId] = s.GetBalance(ctx, config.CURRRPC[currencyId])
	}
}

func (s *WalletService) GetBalance(ctx context.Context, currConfig config.CurrencyRpcConfig) GetBalanceRes {
	var wg sync.WaitGroup
	var res GetBalanceRes = GetBalanceRes{CurrencyConfig: currConfig.Config}

	wg.Add(5)
	go func() { defer wg.Done(); s.SetColdBalanceDetails(ctx, &res) }()
	go func() { defer wg.Done(); s.SetHotBalanceDetails(ctx, currConfig.RpcConfigs, &res) }()
	go func() { defer wg.Done(); s.SetUserBalanceDetails(ctx, &res) }()
	go func() { defer wg.Done(); s.SetPendingWithdraw(ctx, &res) }()
	go func() { defer wg.Done(); s.SetHotLimits(ctx, &res) }()
	wg.Wait()

	s.SetPercent(ctx, &res)

	return res
}

func (s *WalletService) SetColdBalanceDetails(ctx context.Context, res *GetBalanceRes) {
	var (
		symbol   string           = res.CurrencyConfig.Symbol
		cbs      []cb.ColdBalance = s.coldWalletService.GetBalance(ctx, res.CurrencyConfig.Id)
		err      error
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	for _, cb := range cbs {
		var coldBalanceDetail = BalanceDetail{Id: cb.Id, Name: cb.Name, Type: cb.Type}

		coldBalanceDetail.Coin = cb.Balance
		coldBalanceDetail.Address = cb.Address
		coldBalanceDetail.FireblocksName = cb.FireblocksName

		if coldBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(coldBalanceDetail.Coin, symbol); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetColdBalanceDetails)
		}

		if res.TotalColdCoin, err = util.AddCurrency(res.TotalColdCoin, coldBalanceDetail.Coin); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetColdBalanceDetails)
		}

		if res.TotalColdIdr, err = util.AddCurrency(res.TotalColdIdr, coldBalanceDetail.Idr); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetColdBalanceDetails)
		}

		res.ColdBalances = append(res.ColdBalances, coldBalanceDetail)
	}
}

func (s *WalletService) SetHotBalanceDetails(ctx context.Context, rpcConfigs []rc.RpcConfig, res *GetBalanceRes) {
	var (
		symbol   string      = res.CurrencyConfig.Symbol
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	for _, rpcConfig := range rpcConfigs {
		var hotBalanceDetail BalanceDetail = BalanceDetail{Name: rpcConfig.Name, Type: rpcConfig.Type}
		var rpcRes *modulesm.GetBalanceRpcRes

		module, err := s.moduleServices.GetModule(res.CurrencyConfig.Id)
		if err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
			continue
		}

		if rpcRes, err = module.GetBalance(rpcConfig); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
			continue
		}

		hotBalanceDetail.Coin = rpcRes.Balance

		if hotBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(hotBalanceDetail.Coin, symbol); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
		}

		if res.TotalHotCoin, err = util.AddCurrency(res.TotalHotCoin, hotBalanceDetail.Coin); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
		}

		if res.TotalHotIdr, err = util.AddCurrency(res.TotalHotIdr, hotBalanceDetail.Idr); err != nil {
			errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotBalanceDetails)
		}

		res.HotBalances = append(res.HotBalances, hotBalanceDetail)
	}
}

func (s *WalletService) SetUserBalanceDetails(ctx context.Context, res *GetBalanceRes) {
	var (
		tcb                 ub.TotalCoinBalance
		err                 error
		errField            *errs.Error   = nil
		symbol              string        = res.CurrencyConfig.Symbol
		frozenBalanceDetail BalanceDetail = BalanceDetail{Name: "Frozen"}
		liquidBalanceDetail BalanceDetail = BalanceDetail{Name: "Liquid"}
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if tcb, err = s.userBalanceRepo.GetTotalCoinBalance(symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}

	if liquidBalanceDetail.Coin = util.RawToCoin(tcb.Total, 8); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	} else if liquidBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(liquidBalanceDetail.Coin, symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}

	res.UserBalances = append(res.UserBalances, liquidBalanceDetail)

	frozenBalanceDetail.Coin = util.RawToCoin(tcb.TotalFrozen, 8)
	if frozenBalanceDetail.Idr, err = s.marketService.ConvertCoinToIdr(frozenBalanceDetail.Coin, symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}

	res.UserBalances = append(res.UserBalances, frozenBalanceDetail)

	if res.TotalUserCoin, err = util.AddCurrency(liquidBalanceDetail.Coin, frozenBalanceDetail.Coin); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}

	if res.TotalUserIdr, err = util.AddCurrency(liquidBalanceDetail.Idr, frozenBalanceDetail.Idr); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetUserBalanceDetails)
	}
}

func (s *WalletService) SetPendingWithdraw(ctx context.Context, res *GetBalanceRes) {
	var (
		err          error
		errField     *errs.Error = nil
		symbol       string      = res.CurrencyConfig.Symbol
		pendingWDRaw string
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if pendingWDRaw, err = s.withdrawRepo.GetPendingWithdraw(symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPendingWithdraw)
		return
	}

	res.PendingWDCoin = util.RawToCoin(pendingWDRaw, 8)
	if res.PendingWDIdr, err = s.marketService.ConvertCoinToIdr(res.PendingWDCoin, symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPendingWithdraw)
	}
}

func (s *WalletService) SetPercent(ctx context.Context, res *GetBalanceRes) {
	var (
		err      error
		errField *errs.Error = nil
		hotCold  string
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if res.HotPercent, err = util.PercentCurrency(res.TotalHotCoin, res.TotalUserCoin); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPercent)
	}

	if hotCold, err = util.AddCurrency(res.TotalColdCoin, res.TotalHotCoin); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPercent)
	}

	if res.HotColdPercent, err = util.PercentCurrency(hotCold, res.TotalUserCoin); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetPercent)
	}
}

func (s *WalletService) SetHotLimits(ctx context.Context, res *GetBalanceRes) {
	var (
		err      error
		errField *errs.Error = nil
	)

	defer func() {
		if errField != nil {
			logger.ErrorLog(errs.Logged(errField), ctx)
		}
	}()

	if res.HotLimits, err = s.hotLimitRepo.GetBySymbol(res.CurrencyConfig.Symbol); err != nil {
		errField = errs.AssignErr(errs.AddTrace(err), errs.FailedSetHotLimits)
	}
}
