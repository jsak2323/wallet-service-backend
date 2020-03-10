package config

import (
    cc "github.com/btcid/wallet-services/pkg/domain/currencyconfig"
    rc "github.com/btcid/wallet-services/pkg/domain/rpcconfig"
)

type Configuration struct {
    Port string `json:"port"`

    MysqlDbUser     string `json:"mysql_db_user"`
    MysqlDbPass     string `json:"mysql_db_pass"`
    MysqlDbName     string `json:"mysql_db_name"`

    NotificationEmails []string `json:"notification_emails"`

    MailHost            string `json:"mail_host"`
    MailPort            string `json:"mail_port"`
    MailUser            string `json:"mail_user`
    MailEncryptedPass   string `json:"mail_encrypted_pass"`
    MailEncryptionKey   string `json:"mail_encryption_key"`
}

type CurrencyConfigurations struct {
    // BTC CurrencyConfiguration
    ETH CurrencyConfiguration
}

type CurrencyConfiguration struct {
    Config      cc.CurrencyConfig
    RpcConfigs  []rc.RpcConfig
}