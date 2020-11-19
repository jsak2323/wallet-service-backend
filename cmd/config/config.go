package config

import(
    "os"
    "fmt"
    "encoding/json"
    "database/sql"

    _ "github.com/go-sql-driver/mysql"

    mysqldb "github.com/btcid/wallet-services-backend/pkg/database/mysql"
)

var (
    IS_DEV bool
    CONF Configuration
    CURR = make(map[string]CurrencyConfiguration)
)

func init() {
    // todo: changeable from args
    IS_DEV = true
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
    }

    fmt.Println("Done.")
}


