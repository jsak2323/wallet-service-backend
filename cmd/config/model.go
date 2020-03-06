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

    IsCrypto                bool
    IsFinanceEnabled        bool
    IsSingleAddress         bool
    IsUsingMemo             bool
    IsQrCodeEnabled         bool
    IsAddressNoticeEnabled  bool
    IsSeparatedWallet       bool

    QrCodePrefix string

    IsErc20     bool
    IsTrc10     bool
    IsOep4      bool
    IsBep2      bool

    WithdrawFee     float64
    DefaultIdrPrice int

    CmcId       int

    SenderRpcConfig     RpcConfiguration
    ReceiverRpcConfig   RpcConfiguration
    ExtraRpcConfig      RpcConfiguration
}

type RpcConfiguration struct {
    Id          int
    Name        string
    Url         string
    Port        string
    User        string
    Password    string
    Hashkey     string
}