package config

import(
    "fmt"
    "os"
    "encoding/json"
)

var (
    CONF Configuration
    CURR CurrencyConfigurations
)

func init() {
    ReadConfiguration()
    ReadCurrencyConfigurations()
}

func ReadConfiguration() {
    pwd, _ := os.Getwd()
    file, _ := os.Open(pwd+"/cmd/config/config.json")
    defer file.Close()
    decoder := json.NewDecoder(file)
    err := decoder.Decode(&CONF)
    if err != nil {
      fmt.Println("error:", err)
    }
}

func ReadCurrencyConfigurations() {

}