package cron

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"text/template"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	hl "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
	hcw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/util"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
	"github.com/btcid/wallet-services-backend-go/pkg/modules"
)

type CheckBalanceService struct {
	walletService  		*hw.WalletService
	coldWalletService 	*hcw.ColdWalletService
	marketService 		*h.MarketService
	hotLimitRepo   		hl.Repository
	moduleServices 		modules.ModuleServiceMap
}

func NewCheckBalanceService(
	walletService  		*hw.WalletService,
	coldWalletService 	*hcw.ColdWalletService,
	marketService 		*h.MarketService,
	moduleServices 		modules.ModuleServiceMap,
	hotLimitRepo   		hl.Repository,
) *CheckBalanceService {
	return &CheckBalanceService{
		walletService: walletService,
		coldWalletService: coldWalletService,
		marketService: marketService,
		moduleServices: moduleServices,
		hotLimitRepo: hotLimitRepo,
	}
}

func (s *CheckBalanceService) CheckBalanceHandler(w http.ResponseWriter, req *http.Request) {
	for SYMBOL, curr := range config.CURR {	
		var walletBalance hw.WalletBalance = hw.WalletBalance{ 
			TotalColdCoin: "0", TotalColdIdr: "0",
			TotalNodeCoin: "0", TotalNodeIdr: "0",
			TotalUserCoin: "0", TotalUserIdr: "0",
		}
		var wg sync.WaitGroup
		
		wg.Add(3)
		go func() { defer wg.Done(); s.walletService.SetColdBalanceDetails(SYMBOL, &walletBalance) }()
		go func() { defer wg.Done(); s.walletService.SetHotBalanceDetails(SYMBOL, curr.RpcConfigs, &walletBalance) }()
		go func() { defer wg.Done(); s.walletService.SetUserBalanceDetails(SYMBOL, &walletBalance) }()
		wg.Wait()

		// TODO check hot x cold balance, and settle if there's fireblocks available
		s.sendExchangeAlertEmail(SYMBOL, walletBalance)
		s.checkHotLimit(curr.Config, walletBalance)
	}
}

func (s *CheckBalanceService) walletLessThanExchange(walletBalance hw.WalletBalance) bool {
	total, err := util.AddCoin(walletBalance.TotalColdCoin, walletBalance.TotalNodeCoin)
	if err != nil {
		return false
	}

	if compare, err := util.CmpBig(total, walletBalance.TotalUserCoin); err != nil {
		return false
	} else if compare <= 0 { return true }

	return false
}

func (s *CheckBalanceService) sendExchangeAlertEmail(symbol string, walletBalance hw.WalletBalance) {
	logger.Log(" - CheckBalanceService -- Sending balance report email ...")

    subject := "Balance Alert: "+symbol
    message := `Content-Type: text/html; charset=UTF-8 \r\n`

	buf := &bytes.Buffer{}
	
	t, err := template.ParseFiles("views/email/exchange_alert.html")
	if err != nil { logger.ErrorLog(err.Error()) }

	fmt.Println(walletBalance)
	err = t.Execute(buf, struct {WalletBalance hw.WalletBalance}{WalletBalance: walletBalance})

	message = message + buf.String()

	fmt.Println(message)

	recipients := []string{} // TODO user email with certain role

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog(err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

func (s *CheckBalanceService) sendReportEmail(symbol string, walletBalance hw.WalletBalance) {
	// TODO create email body
	logger.Log(" - CheckBalanceService -- Sending balance report email ...")

    subject := "Balance Report: "+symbol
    message := ""
    recipients := []string{} // TODO user email with certain role

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog(err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

func (s *CheckBalanceService) sendHotLimitAlertEmail(symbol string, walletBalance hw.WalletBalance, limits map[string]hl.HotLimit) {
	logger.Log(" - CheckBalanceService -- Sending hot limit report email ...")

    subject := "Hot Limit Alert: "+symbol
    message := `Content-Type: text/html; charset=UTF-8 \r\n`

	buf := &bytes.Buffer{}
	
	t, err := template.ParseFiles("views/email/hot_limit_alert.html")
	if err != nil { logger.ErrorLog(err.Error()) }

	err = t.Execute(buf, struct {
			Symbol string
			WalletBalance hw.WalletBalance
			Limits map[string]hl.HotLimit
		}{
			Symbol: symbol,
			WalletBalance: walletBalance,
			Limits: limits,
		})

	message = message + buf.String()

	fmt.Println(message)

	recipients := []string{} // TODO user email with certain role

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog(err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

func (s *CheckBalanceService) checkHotLimit(currency cc.CurrencyConfig, walletBalance hw.WalletBalance) {
	limits, err := s.hotLimitRepo.GetByCurrencyId(currency.Id)
	if err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetByCurrencyId err: "+err.Error())
		return
	}

	senderRpc, err := rc.GetSenderFromList(config.CURR[currency.Symbol].RpcConfigs)
	if err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetSenderFromList err: "+err.Error())
		return
	}

	// check if hot storage is greater than bottom soft limit
	if compare, err := util.CmpBig(walletBalance.TotalNodeIdr, limits[hl.TopSoftType].Amount); err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") CmpBig err: "+err.Error())
		return
	} else if compare == 1 {
		amount, err := util.SubIdr(walletBalance.TotalNodeIdr, limits[hl.TargetType].Amount)
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") SubIdr err: "+err.Error())
			return
		}

		amount, err = s.marketService.ConvertIdrToCoin(amount, currency.Symbol)
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") ConvertIdrToCoin err: "+err.Error())
			return
		}

		address, err := s.coldWalletService.GetDepositAddress(currency.Id)
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetDepositAddress err: "+err.Error())
			return
		}

		// TODO get memo
		memo := ""

		// TODO print attempt log
		if res, err := s.moduleServices[currency.Symbol].SendToAddress(senderRpc, amount, address, memo); err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") SendToAddress err: "+err.Error())
			return
		} else {
			logger.InfoLog("checkHotLimit("+currency.Symbol+") SendToAddress res: "+res.TxHash, &http.Request{})
			return
		}
	}

	// check if hot storage is less than bottom soft limit
	if compare, err := util.CmpBig(walletBalance.TotalNodeIdr, limits[hl.BottomSoftType].Amount); err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") CmpBig err: "+err.Error())
		return
	} else if compare == -1 || true {
		amount, err := util.SubIdr(limits[hl.TargetType].Amount, walletBalance.TotalNodeIdr)
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") SubIdr err: "+err.Error())
			return
		}

		amount, err = s.marketService.ConvertIdrToCoin(amount, currency.Symbol)
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") ConvertIdrToCoin err: "+err.Error())
			return
		}

		if currency.FireblocksName == "" {
			// TODO print attempt log
			s.sendHotLimitAlertEmail(currency.Symbol, walletBalance, limits)
		} else {
			// TODO print attempt log
			logger.InfoLog("Sending from fireblocks cold to hot...", &http.Request{})
			res, err := fireblocks.CreateTransaction(fireblocks.CreateTransactionReq{
				AssetId: currency.FireblocksName,
				Amount: amount,
				Source: fireblocks.TransactionAccount{Type: fireblocks.VaultAccountType, Id: config.CONF.FireblocksColdVaultId},
				Destination: fireblocks.TransactionAccount{Type: fireblocks.InternalWalletType, Id: config.CONF.FireblocksHotVaultId},
			})
			if err != nil {
				logger.ErrorLog(" - SendToHotHandler fireblocks.CreateTransaction err: " + err.Error())
				return
			}
		
			if res.Error != "" {
				logger.ErrorLog(" - SendToHotHandler fireblocks.CreateTransaction err: " + res.Error)
				return
			}
		}
		// TODO print success log
	}
}