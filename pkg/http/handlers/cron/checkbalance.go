package cron

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/btcid/wallet-services-backend-go/cmd/config"
	cb "github.com/btcid/wallet-services-backend-go/pkg/domain/coldbalance"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	hl "github.com/btcid/wallet-services-backend-go/pkg/domain/hotlimit"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	h "github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	hw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet"
	hcw "github.com/btcid/wallet-services-backend-go/pkg/http/handlers/wallet/cold"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/fireblocks"
	"github.com/btcid/wallet-services-backend-go/pkg/lib/telegram"
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
	var walletBalances []hw.GetBalanceRes

	for _, curr := range config.CURR {	
		walletBalance := s.walletService.GetBalance(curr)
		
		s.checkUserBalance(walletBalance)
		// s.checkHotLimit(curr.Config, walletBalance)

		walletBalances = append(walletBalances, walletBalance)
	}

	s.sendReportEmail(walletBalances)
}

func (s *CheckBalanceService) checkUserBalance(walletBalance hw.GetBalanceRes) {
	var err error
	var total string = "0"
	var cmpResult int

	total, err = util.AddCoin(total, walletBalance.TotalColdCoin)
	total, err = util.AddCoin(total, walletBalance.TotalHotCoin)

	if cmpResult, err = util.CmpBig(total, walletBalance.TotalUserCoin); err != nil {
		// TODO log error
	} else if cmpResult == -1 {
		s.walletService.FormatWalletBalanceCurrency(&walletBalance)
		s.sendUserBalanceAlertTelegram(walletBalance, total)
		s.sendUserBalanceAlertEmail(walletBalance, total)
	}
}

func (s *CheckBalanceService) sendUserBalanceAlertTelegram(walletBalance hw.GetBalanceRes, total string) {
	logger.Log(" - CheckBalanceService -- Sending user balance alert telegram ...")

	var symbol string = walletBalance.CurrencyConfig.Symbol
	var unit string = walletBalance.CurrencyConfig.Unit
	var sb strings.Builder

	sb.WriteString("User Balance Alert: "+ symbol + "\n")
	sb.WriteString("Total Hot: " + walletBalance.TotalHotCoin + " " + unit + "\n")
	sb.WriteString("Total Cold: " + walletBalance.TotalColdCoin + " " + unit + "\n")
	sb.WriteString("Total Wallet: " + total + " " + unit + "\n")
	sb.WriteString("User Balance: " + walletBalance.TotalUserCoin + " " + unit)
	
	telegram.SendMessage(sb.String())
}

func (s *CheckBalanceService) sendUserBalanceAlertEmail(walletBalance hw.GetBalanceRes, total string) {
	logger.Log(" - CheckBalanceService -- Sending user balance alert email ...")

	subject := "User Balance Alert: "+walletBalance.CurrencyConfig.Symbol
    message := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buf := &bytes.Buffer{}
	
	t, err := template.ParseFiles("views/email/user_balance_alert.html")
	if err != nil { logger.ErrorLog(err.Error()) }

	err = t.Execute(buf, struct {
			WalletBalance hw.GetBalanceRes
			Total string
		}{
			WalletBalance: walletBalance,
			Total: total,
		})

	message = message + buf.String()

	recipients := config.CONF.NotificationEmails // TODO user email with certain role

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog(err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

func (s *CheckBalanceService) sendReportEmail(walletBalances []hw.GetBalanceRes) {
	logger.Log(" - CheckBalanceService -- Sending balance report email ...")

    subject := "Balance Report "
    message := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buf := &bytes.Buffer{}
	
	t, err := template.ParseFiles("views/email/report.html")
	if err != nil { logger.ErrorLog(err.Error()) }

	for i := range walletBalances {
		s.walletService.FormatWalletBalanceCurrency(&walletBalances[i])
	}

	err = t.Execute(buf, struct {WalletBalances []hw.GetBalanceRes}{WalletBalances: walletBalances})

	message = message + buf.String()

    recipients := config.CONF.NotificationEmails // TODO user email with certain role

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog(err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

func (s *CheckBalanceService) sendHotLimitAlertEmail(symbol string, walletBalance hw.GetBalanceRes, limits hl.HotLimit) {
	logger.Log(" - CheckBalanceService -- Sending hot limit report email ...")

    subject := "Hot Limit Alert: "+symbol
    message := `Content-Type: text/html; charset=UTF-8 \r\n`

	buf := &bytes.Buffer{}
	
	t, err := template.ParseFiles("views/email/hot_limit_alert.html")
	if err != nil { logger.ErrorLog(err.Error()) }

	err = t.Execute(buf, struct {
			Symbol string
			WalletBalance hw.GetBalanceRes
			Limits hl.HotLimit
		}{
			Symbol: symbol,
			WalletBalance: walletBalance,
			Limits: limits,
		})

	message = message + buf.String()

	recipients := config.CONF.NotificationEmails // TODO user email with certain role

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog(err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

func (s *CheckBalanceService) checkHotLimit(currency cc.CurrencyConfig, walletBalance hw.GetBalanceRes) {
	limits, err := s.hotLimitRepo.GetBySymbol(currency.Symbol)
	if err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetByCurrencyId err: "+err.Error())
		return
	}

	senderRpc, err := rc.GetSenderFromList(config.CURR[currency.Symbol].RpcConfigs)
	if err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetSenderFromList err: "+err.Error())
		return
	}

	coldWallet, err := s.coldWalletService.SettlementWallet(currency.Id)
	if err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetDepositAddress err: "+err.Error())
		return
	}

	// check if hot storage is greater than bottom soft limit
	if compare, err := util.CmpBig(walletBalance.TotalHotIdr, limits[hl.TopSoftType]); err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") CmpBig err: "+err.Error())
		return
	} else if compare == 1 {
		amount, err := util.SubIdr(walletBalance.TotalHotIdr, limits[hl.TargetType])
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") SubIdr err: "+err.Error())
			return
		}

		amount, err = s.marketService.ConvertIdrToCoin(amount, currency.Symbol)
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") ConvertIdrToCoin err: "+err.Error())
			return
		}

		address := coldWallet.Address
		// TODO get memo
		memo := ""

		// TODO print attempt log
		if res, err := s.moduleServices[currency.Symbol].SendToAddress(senderRpc, amount, address, memo); err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") SendToAddress err: "+err.Error())
			return
		} else {
			// TODO update db
			logger.InfoLog("checkHotLimit("+currency.Symbol+") SendToAddress res: "+res.TxHash, &http.Request{})
			return
		}
	}

	// check if hot storage is less than bottom soft limit
	if compare, err := util.CmpBig(walletBalance.TotalHotIdr, limits[hl.BottomSoftType]); err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") CmpBig err: "+err.Error())
		return
	} else if compare == -1 {
		amount, err := util.SubIdr(limits[hl.TargetType], walletBalance.TotalHotIdr)
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") SubIdr err: "+err.Error())
			return
		}

		amount, err = s.marketService.ConvertIdrToCoin(amount, currency.Symbol)
		if err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") ConvertIdrToCoin err: "+err.Error())
			return
		}

		// TODO check zero amount
		if amount == "0" || amount == "" { logger.ErrorLog("checkHotLimit("+currency.Symbol+") zero amount "); return }

		if coldWallet.Type == cb.ColdType {
			logger.InfoLog("Sending hot limit alert for non fireblocks...", &http.Request{})
			s.sendHotLimitAlertEmail(currency.Symbol, walletBalance, limits)
		} else if coldWallet.Type == cb.FbWarmType || coldWallet.Type == cb.FbColdType {
			// TODO print attempt log
			logger.InfoLog("Sending from fireblocks cold to hot...", &http.Request{})
			if res, err := fireblocks.CreateTransaction(fireblocks.CreateTransactionReq{
				AssetId: coldWallet.FireblocksName,
				Amount: amount,
				Source: fireblocks.TransactionAccount{
					Type: fireblocks.VaultAccountType, 
					Id: hcw.FireblocksVaultAccountId(coldWallet.Type),
				},
				Destination: fireblocks.TransactionAccount{
					Type: fireblocks.InternalWalletType, 
					Id: config.CONF.FireblocksHotVaultId,
				},
			}); err != nil {
				logger.ErrorLog(" - SendToHotHandler fireblocks.CreateTransaction err: " + err.Error())
				return
			} else if res.Error != "" {
				logger.ErrorLog(" - SendToHotHandler fireblocks.CreateTransaction err: " + res.Error)
				return
			} else { logger.InfoLog("checkHotLimit("+currency.Symbol+") Sent from fireblocks res: "+res.Id, &http.Request{}) }
		}		
	}
}