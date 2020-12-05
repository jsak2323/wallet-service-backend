package config

import(
    "os"
    "fmt"
    "database/sql"
    "encoding/json"

    _ "github.com/go-sql-driver/mysql"

    mysqldb "github.com/btcid/wallet-services-backend/pkg/database/mysql"
    cc "github.com/btcid/wallet-services-backend/pkg/domain/currencyconfig"
    rc "github.com/btcid/wallet-services-backend/pkg/domain/rpcconfig"
)

var (
    IS_DEV bool
    CONF Configuration
    CURR = make(map[string]CurrencyConfiguration)
    SYMBOLS = make(map[int]string)
)

type Configuration struct {
    Port string `json:"port"`

    MysqlDbUser     string `json:"mysql_db_user"`
    MysqlDbPass     string `json:"mysql_db_pass"`
    MysqlDbName     string `json:"mysql_db_name"`

    NotificationEmails []string `json:"notification_emails"`

    AuthorizedIps []string `json:"authorized_ips"`

    MailHost            string `json:"mail_host"`
    MailPort            string `json:"mail_port"`
    MailUser            string `json:"mail_user"`
    MailEncryptedPass   string `json:"mail_encrypted_pass"`
    MailEncryptionKey   string `json:"mail_encryption_key"`

    CryptoApisKey       string `json:"crypto_apis_key"`
    InfuraProjectId     string `json:"infura_project_id"`
}

type CurrencyConfiguration struct {
    Config      cc.CurrencyConfig
    RpcConfigs  []rc.RpcConfig
}

func init() {
    IS_DEV = os.Getenv("PRODUCTION") != "true"

    fmt.Println()
    env := "development"
    if !IS_DEV { env = "production" }
    fmt.Println("Environment: "+env)
    
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
    file, _ := os.Open(gopath+"/src/github.com/btcid/wallet-services-backend/cmd/config/json/"+configFilename)
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
    dbUser   := CONF.MysqlDbUser
    dbPass   := CONF.MysqlDbPass
    dbName   := CONF.MysqlDbName

    db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
    if err != nil { panic(err.Error()) }
    return db
}

func LoadCurrencyConfigs() {
    fmt.Print("Loading Currency Configurations ... ")

    mysqlDb := MysqlDbConn()
    defer mysqlDb.Close()

    currencyConfigRepo  := mysqldb.NewMysqlCurrencyConfigRepository(mysqlDb)
    rpcConfigRepo       := mysqldb.NewMysqlRpcConfigRepository(mysqlDb)

    currencyConfigs, err := currencyConfigRepo.GetAll()
    if err != nil { 
        panic("Unexpected error when executing currencyConfigRepo.GetAll(). Error: "+err.Error()) 
    }

    for _, currencyConfig := range currencyConfigs {
        rpcConfigs, err := rpcConfigRepo.GetByCurrencyId(currencyConfig.Id)
        if err != nil {
            panic("Unexpected error when executing rpcConfigRepo.GetByCurrencyId(currencyConfig.Id). Error: "+err.Error())
        }

        CURR[currencyConfig.Symbol] = CurrencyConfiguration{
            Config      : currencyConfig,
            RpcConfigs  : rpcConfigs,
        }

        SYMBOLS[currencyConfig.Id] = currencyConfig.Symbol
    }

    fmt.Println("Done.")
}


