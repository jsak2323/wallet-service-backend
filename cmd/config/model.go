package config

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
    Id              int
    Symbol          string
    Name            string
    NameUppercase   string
    NameLowercase   string
    Unit            string
    TokenType       string

    IsFinanceEnabled        bool
    IsSingleAddress         bool
    IsUsingMemo             bool
    IsQrCodeEnabled         bool
    IsAddressNoticeEnabled  bool

    QrCodePrefix string

    WithdrawFee     float64
    DefaultIdrPrice int

    CmcId       int
    RpcConfigs  []RpcConfiguration
}

type RpcConfiguration struct {
    Id          int
    CurrencyId  int
    Type        string
    Name        string
    Host        string
    Port        string
    Path        string
    User        string
    Password    string
    Hashkey     string
}