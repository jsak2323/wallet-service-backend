package currencyconfig

type CurrencyConfig struct {
    Id                      int
    Symbol                  string
    Name                    string
    NameUppercase           string
    NameLowercase           string
    Unit                    string
    TokenType               string

    IsFinanceEnabled        bool
    IsSingleAddress         bool
    IsUsingMemo             bool
    IsQrCodeEnabled         bool
    IsAddressNoticeEnabled  bool

    QrCodePrefix            string
    WithdrawFee             string
    HealthyBlockDiff        int
    DefaultIdrPrice         int
    CmcId                   int
    LastUpdated             string
}


