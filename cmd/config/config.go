package config

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"

	mysqldb "github.com/btcid/wallet-services-backend-go/pkg/database/mysql"
	cc "github.com/btcid/wallet-services-backend-go/pkg/domain/currencyconfig"
	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
)

var (
	IS_DEV           bool
	CONF             Configuration
	ErrorMailCount   int
	FirstHealthCheck bool

	CURRRPC = make(map[int]CurrencyRpcConfig)
	CURRID  = make(map[string]map[string]int) // currency_id by symbol and token_typen
	SYMBOLS = make(map[int]string)
)

type Configuration struct {
	Port string `json:"port"`

	MysqlDbUser string `json:"mysql_db_user"`
	MysqlDbPass string `json:"mysql_db_pass"`
	MysqlDbName string `json:"mysql_db_name"`

	ExchangeSlaveMysqlDbHost string `json:"exchange_slave_mysql_db_host"`
	ExchangeSlaveMysqlDbUser string `json:"exchange_slave_mysql_db_user"`
	ExchangeSlaveMysqlDbPass string `json:"exchange_slave_mysql_db_pass"`
	ExchangeSlaveMysqlDbName string `json:"exchange_slave_mysql_db_name"`

	ExchangeHost string `json:"exchange_host"`

	NotificationEmails []string `json:"notification_emails"`

	AuthorizedIps []string `json:"authorized_ips"`

	MailHost          string `json:"mail_host"`
	MailPort          string `json:"mail_port"`
	MailUser          string `json:"mail_user"`
	MailAddress       string `json:"mail_address"`
	MailEncryptedPass string `json:"mail_encrypted_pass"`
	MailEncryptionKey string `json:"mail_encryption_key"`

	SessionErrorMailNotifLimit int `json:"session_error_mail_notif_limit"`

	CryptoApisKey   string `json:"crypto_apis_key"`
	InfuraProjectId string `json:"infura_project_id"`

	JWTSecret string `json:"jwt_secret"`

	FireblocksHost          string `json:"fireblocks_host"`
	FireblocksServeruser    string `json:"fireblocks_serveruser"`
	FireblocksServerpass    string `json:"fireblocks_serverpass"`
	FireblocksServerhashkey string `json:"fireblocks_serverhashkey"`
	FireblocksColdVaultId   string `json:"fireblocks_cold_vault_id"`
	FireblocksWarmVaultId   string `json:"fireblocks_warm_vault_id"`
	FireblocksHotVaultId    string `json:"fireblocks_hot_vault_id"`

	EthEncryptKeyEncrypted string `json:"eth_encrypt_key_encrypted"`
	EthEncryptKeyKey       string `json:"eth_encrypt_key_key"`

	TelegramBotToken    string `json:"telegram_bot_token"`
	TelegramAlertChatId string `json:"telegram_alert_chat_id"`
}

type CurrencyRpcConfig struct {
	Config     cc.CurrencyConfig
	RpcConfigs []rc.RpcConfig
}

func init() {
	IS_DEV = os.Getenv("PRODUCTION") != "true"

	fmt.Println()
	env := "development"
	if !IS_DEV {
		env = "production"
	}
	fmt.Println("Environment: " + env)

	LoadAppConfig()
	LoadCurrencyConfigs()
}

func LoadAppConfig() {
	configFilename := "config.json"
	if IS_DEV {
		configFilename = "config-dev.json"
	}

	fmt.Print("Loading App Configuration ... ")
	gopath := os.Getenv("GOPATH")
	file, _ := os.Open(gopath + "/src/github.com/btcid/wallet-services-backend-go/cmd/config/json/" + configFilename)
	defer file.Close()
	decoder := json.NewDecoder(file)
	err := decoder.Decode(&CONF)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println("Done.")
}

func MysqlDbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := CONF.MysqlDbUser
	dbPass := CONF.MysqlDbPass
	dbName := CONF.MysqlDbName

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	// preventive idle connection error
	db.SetMaxIdleConns(0)
	return db
}

func ExchangeSlaveMysqlDbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbHost := CONF.ExchangeSlaveMysqlDbHost
	dbUser := CONF.ExchangeSlaveMysqlDbUser
	dbPass := CONF.ExchangeSlaveMysqlDbPass
	dbName := CONF.ExchangeSlaveMysqlDbName

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHost+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func LoadCurrencyConfigs() {
	fmt.Print("Loading Currency Configurations ... ")

	mysqlDb := MysqlDbConn()
	defer mysqlDb.Close()

	currencyConfigRepo := mysqldb.NewMysqlCurrencyConfigRepository(mysqlDb)
	rpcConfigRepo := mysqldb.NewMysqlRpcConfigRepository(mysqlDb)

	currencyConfigs, err := currencyConfigRepo.GetAll()
	if err != nil {
		panic("Unexpected error when executing currencyConfigRepo.GetAll(). Error: " + err.Error())
	}

	for _, currencyConfig := range currencyConfigs {
		getSymbol := currencyConfig.Symbol
		if currencyConfig.ParentSymbol != "" {
			getSymbol = currencyConfig.ParentSymbol
		}

		rpcConfigs, err := rpcConfigRepo.GetByCurrencyId(currencyConfig.Id)
		if err != nil {
			panic("Unexpected error when executing rpcConfigRepo.GetByCurrencySymbol(getSymbol). " + getSymbol + " Error: " + err.Error())
		}

		CURRRPC[currencyConfig.Id] = CurrencyRpcConfig{
			Config:     currencyConfig,
			RpcConfigs: rpcConfigs,
		}

		_, ok := CURRID[currencyConfig.Symbol]
		if !ok {
			CURRID[currencyConfig.Symbol] = make(map[string]int)
		}
		CURRID[currencyConfig.Symbol][currencyConfig.TokenType] = currencyConfig.Id

		SYMBOLS[currencyConfig.Id] = currencyConfig.Symbol
	}

	fmt.Println("Done.")
}

func GetCurrencyBySymbol(symbol string) (result []cc.CurrencyConfig, err error) {
	tokenTypes, ok := CURRID[symbol]
	if !ok {
		return []cc.CurrencyConfig{}, errs.AddTrace(errors.New("symbol not found " + symbol))
	}

	for _, currencyConfigId := range tokenTypes {
		result = append(result, CURRRPC[currencyConfigId].Config)
	}

	return result, nil
}

func GetCurrencyBySymbolTokenType(symbol, tokenType string) (cc.CurrencyConfig, error) {
	tokenTypes, ok := CURRID[symbol]
	if !ok {
		return cc.CurrencyConfig{}, errs.AddTrace(errors.New("symbol not found  " + symbol))
	}

	currencyId, ok := tokenTypes[tokenType]
	if !ok {
		return cc.CurrencyConfig{}, errs.AddTrace(errors.New("token_type not found " + tokenType))
	}

	return CURRRPC[currencyId].Config, nil
}

func GetRpcConfigByType(currencyConfigId int, rpcConfigType string) (rc.RpcConfig, error) {
	for _, rpcConfig := range CURRRPC[currencyConfigId].RpcConfigs {
		if rpcConfig.Type == rpcConfigType || rpcConfig.Type == "master" {
			return rpcConfig, nil
		}
	}
	return rc.RpcConfig{}, errs.AddTrace(errors.New("RpcConfig not found."))
}
