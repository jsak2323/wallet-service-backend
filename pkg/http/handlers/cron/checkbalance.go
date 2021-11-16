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
	"github.com/btcid/wallet-services-backend-go/pkg/domain/user"
)

type CheckBalanceService struct {
	walletService  		*hw.WalletService
	coldWalletService 	*hcw.ColdWalletService
	marketService 		*h.MarketService
	moduleServices 		*modules.ModuleServiceMap
	hotLimitRepo   		hl.Repository
	userRepo			user.Repository
}

func NewCheckBalanceService(
	walletService  		*hw.WalletService,
	coldWalletService 	*hcw.ColdWalletService,
	marketService 		*h.MarketService,
	moduleServices 		*modules.ModuleServiceMap,
	hotLimitRepo   		hl.Repository,
	userRepo 			user.Repository,
) *CheckBalanceService {
	return &CheckBalanceService{
		walletService: walletService,
		coldWalletService: coldWalletService,
		marketService: marketService,
		moduleServices: moduleServices,
		hotLimitRepo: hotLimitRepo,
		userRepo: userRepo,
	}
}

const adminRoleName = "admin"
const defaultMemo = "0000"

func (s *CheckBalanceService) CheckBalanceHandler() {
	var walletBalances []hw.GetBalanceRes

	for _, curr := range config.CURRRPC {	
		walletBalance := s.walletService.GetBalance(curr)
	
		s.checkUserBalance(walletBalance)
		s.checkHotLimit(curr.Config, walletBalance)

		walletBalances = append(walletBalances, walletBalance)
	}

	s.sendReportEmail(walletBalances)
}

func (s *CheckBalanceService) checkUserBalance(walletBalance hw.GetBalanceRes) {
	logger.Log(" - CheckBalanceService -- Checking "+walletBalance.CurrencyConfig.Symbol+" "+walletBalance.CurrencyConfig.TokenType+" user balance...")

	var err error
	var totalCoin string = "0"
	var cmpResult int

	totalCoin, err = util.AddCoin(totalCoin, walletBalance.TotalColdCoin)
	if err != nil { logger.ErrorLog("checkUserBalance("+walletBalance.CurrencyConfig.Symbol+") AddCoin err: "+err.Error()) }
	
	totalCoin, err = util.AddCoin(totalCoin, walletBalance.TotalHotCoin)
	if err != nil { logger.ErrorLog("checkUserBalance("+walletBalance.CurrencyConfig.Symbol+") AddCoin err: "+err.Error()) }

	if cmpResult, err = util.CmpBig(totalCoin, walletBalance.TotalUserCoin); err != nil {
		if err != nil { logger.ErrorLog("checkUserBalance("+walletBalance.CurrencyConfig.Symbol+") CmpBig err: "+err.Error()) }
	} else if cmpResult == -1 {
		walletBalanceFormatted := s.walletService.FormatWalletBalanceCurrency(walletBalance)
		s.sendUserBalanceAlertTelegram(walletBalanceFormatted, totalCoin)
		s.sendUserBalanceAlertEmail(walletBalanceFormatted, totalCoin)
	} else { logger.Log(" - CheckBalanceService -- Finished checking "+walletBalance.CurrencyConfig.Symbol+" "+walletBalance.CurrencyConfig.TokenType+" user balance") }
}

func (s *CheckBalanceService) sendUserBalanceAlertTelegram(walletBalance hw.GetBalanceRes, totalCoin string) {
	logger.Log(" - CheckBalanceService -- Sending user balance alert telegram ...")

	var symbol string = walletBalance.CurrencyConfig.Symbol
	var unit string = walletBalance.CurrencyConfig.Unit
	var sb strings.Builder

	sb.WriteString("User Balance Alert: "+ symbol + "\n")
	sb.WriteString("Total Hot: " + walletBalance.TotalHotCoin + " " + unit + "\n")
	sb.WriteString("Total Cold: " + walletBalance.TotalColdCoin + " " + unit + "\n")
	sb.WriteString("Total Wallet: " + totalCoin + " " + unit + "\n")
	sb.WriteString("User Balance: " + walletBalance.TotalUserCoin + " " + unit)
	
	telegram.SendMessage(sb.String())
}

func (s *CheckBalanceService) sendUserBalanceAlertEmail(walletBalance hw.GetBalanceRes, totalCoin string) {
	logger.Log(" - CheckBalanceService -- Sending user balance alert email ...")

	subject := "User Balance Alert: "+walletBalance.CurrencyConfig.Symbol
    message := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buf := &bytes.Buffer{}
	
	t, err := template.ParseFiles("views/email/user_balance_alert.html")
	if err != nil { logger.ErrorLog("sendUserBalanceAlertEmail("+walletBalance.CurrencyConfig.Symbol+") ParseFiles err: "+err.Error()) }

	err = t.Execute(buf, struct {
			WalletBalance hw.GetBalanceRes
			TotalCoin string
		}{
			WalletBalance: walletBalance,
			TotalCoin: totalCoin,
		})

	message = message + buf.String()

	recipients, err := s.userRepo.GetEmailsByRole(adminRoleName)
	if err != nil { logger.ErrorLog("sendUserBalanceAlertEmail("+walletBalance.CurrencyConfig.Symbol+") GetEmailsByRole err: "+err.Error()) }
	
	isEmailSent, err := util.SendEmail(subject, message, recipients)
	if err != nil { logger.ErrorLog("sendUserBalanceAlertEmail("+walletBalance.CurrencyConfig.Symbol+") SendEmail err: "+err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

func (s *CheckBalanceService) sendReportEmail(walletBalances []hw.GetBalanceRes) {
	logger.Log(" - CheckBalanceService -- Sending balance report email ...")

    subject := "Balance Report "
    message := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buf := &bytes.Buffer{}
	
	t, err := template.ParseFiles("views/email/report.html")
    if err != nil { logger.ErrorLog("sendReportEmail err: "+err.Error()) }
	
	walletBalancesFormatted := []hw.GetBalanceRes{}
	for i := range walletBalances {
		walletBalanceFormatted := s.walletService.FormatWalletBalanceCurrency(walletBalances[i])
		walletBalancesFormatted = append(walletBalancesFormatted, walletBalanceFormatted)
	}

	err = t.Execute(buf, struct {WalletBalances []hw.GetBalanceRes}{WalletBalances: walletBalancesFormatted})

	message = message + buf.String()

    recipients, err := s.userRepo.GetEmailsByRole(adminRoleName)
	if err != nil { logger.ErrorLog("sendReportEmail GetEmailsByRole err: "+err.Error()) }

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog("sendReportEmail err: "+err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

func (s *CheckBalanceService) sendHotLimitAlertEmail(symbol string, walletBalance hw.GetBalanceRes, limits hl.HotLimit) {
	logger.Log(" - CheckBalanceService -- Sending "+symbol+" hot limit alert for non fireblocks...")
	
    subject := "Hot Limit Alert: "+symbol
    message := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	buf := &bytes.Buffer{}
	
	t, err := template.ParseFiles("views/email/hot_limit_alert.html")
	if err != nil { logger.ErrorLog("sendHotLimitAlertEmail("+symbol+") ParseFiles err: "+err.Error()) }

	limitsFormatted := make(map[string]string, len(limits))
	for key, limit := range limits {
		limitsFormatted[key] = util.FormatCurrency(limit) + " IDR"
	}

	err = t.Execute(buf, struct {
			Symbol string
			WalletBalance hw.GetBalanceRes
			Limits hl.HotLimit
		}{
			Symbol: symbol,
			WalletBalance: s.walletService.FormatWalletBalanceCurrency(walletBalance),
			Limits: limitsFormatted,
		})

	message = message + buf.String()

	recipients, err := s.userRepo.GetEmailsByRole(adminRoleName)
	if err != nil { logger.ErrorLog("sendHotLimitAlertEmail("+walletBalance.CurrencyConfig.Symbol+") GetEmailsByRole err: "+err.Error()) }

    isEmailSent, err := util.SendEmail(subject, message, recipients)
    if err != nil { logger.ErrorLog("sendHotLimitAlertEmail("+symbol+") SendEmail err: "+err.Error()) }
    logger.Log(" - CheckBalanceService -- Is balance report email sent: "+strconv.FormatBool(isEmailSent))
}

// TODO check per network
func (s *CheckBalanceService) checkHotLimit(currency cc.CurrencyConfig, walletBalance hw.GetBalanceRes) {
	logger.Log(" - CheckBalanceService -- Checking "+currency.Symbol+" " + currency.TokenType + " hot limit...")
	
	limits, err := s.hotLimitRepo.GetBySymbol(currency.Symbol)
	if err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetByCurrencyId err: "+err.Error())
		return
	}

	senderRpc, err := util.GetRpcConfigByType(currency.Id, rc.SenderRpcType)
	if err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetSenderFromList err: "+err.Error())
		return
	}

	coldWallet, err := s.coldWalletService.SettlementWallet(currency.Id)
	if err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") GetDepositAddress err: "+err.Error())
		return
	}

	// check if hot storage is greater than top soft limit
	if compare, err := util.CmpBig(walletBalance.TotalHotIdr, limits[hl.TopSoftType]); err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") CmpBig err: "+err.Error())
		return
	} else if compare == 1 {
		logger.Log(" - CheckBalanceService -- Hot balance "+currency.Symbol+" is greater than top soft limit")

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
		memo := defaultMemo

		module, err := s.moduleServices.GetModule(currency.Id)
		if err != nil {
			logger.ErrorLog("s.moduleServices.GetModule err: "+err.Error())
			return
		}
		
		logger.Log(" - CheckBalanceService -- Sending "+currency.Symbol+" from hot to cold...")
		if res, err := module.SendToAddress(senderRpc, amount, address, memo); err != nil {
			logger.ErrorLog("checkHotLimit("+currency.Symbol+") SendToAddress err: "+err.Error())
			return
		} else {
			logger.Log(" - CheckBalanceService -- checkHotLimit("+currency.Symbol+") SendToAddress sent with tx: "+res.TxHash)
			balanceToUpdate, err := util.AddCoin(coldWallet.Balance, amount)
			if err != nil {
				logger.ErrorLog("checkHotLimit("+currency.Symbol+") AddCoin err: "+err.Error())
			}

			if err = s.coldWalletService.UpdateBalance(coldWallet.Id, balanceToUpdate); err != nil {

			}

			logger.Log(" - CheckBalanceService -- checkHotLimit("+currency.Symbol+") UpdateBalance updated with amount: "+balanceToUpdate)
			return
		}
	}

	// check if hot storage is less than bottom soft limit
	if compare, err := util.CmpBig(walletBalance.TotalHotIdr, limits[hl.BottomSoftType]); err != nil {
		logger.ErrorLog("checkHotLimit("+currency.Symbol+") CmpBig err: "+err.Error())
		return
	} else if compare == -1 {
		logger.InfoLog("Hot balance "+currency.Symbol+" is less than bottom soft limit", &http.Request{})

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

		if amount == "0" || amount == "" { logger.ErrorLog("checkHotLimit("+currency.Symbol+") zero amount "); return }

		if coldWallet.Type == cb.ColdType {
			s.sendHotLimitAlertEmail(currency.Symbol, walletBalance, limits)
		} else if coldWallet.Type == cb.FbWarmType || coldWallet.Type == cb.FbColdType {
			logger.InfoLog("Sending "+currency.Symbol+" from fireblocks cold to hot...", &http.Request{})
			
			vaultAccountId, err := hcw.FireblocksVaultAccountId(coldWallet.Type)
			if err != nil { logger.ErrorLog(" - checkHotLimit fireblocks.FireblocksVaultAccountId err: " + err.Error()) }

			if res, err := fireblocks.CreateTransaction(fireblocks.CreateTransactionReq{
				AssetId: coldWallet.FireblocksName,
				Amount: amount,
				Source: fireblocks.TransactionAccount{
					Type: fireblocks.VaultAccountType, 
					Id: vaultAccountId,
				},
				Destination: fireblocks.TransactionAccount{
					Type: fireblocks.InternalWalletType, 
					Id: config.CONF.FireblocksHotVaultId,
				},
			}); err != nil {
				logger.ErrorLog(" - checkHotLimit fireblocks.CreateTransaction err: " + err.Error())
				return
			} else if res.Error != "" {
				logger.ErrorLog(" - checkHotLimit fireblocks.CreateTransaction err: " + res.Error)
				return
			} else { logger.InfoLog("checkHotLimit("+currency.Symbol+") Sent from fireblocks res: "+res.Id, &http.Request{}) }
		}		
	}

	logger.Log(" - CheckBalanceService -- Finished checking "+walletBalance.CurrencyConfig.Symbol+" "+walletBalance.CurrencyConfig.TokenType+" hot limit")
}