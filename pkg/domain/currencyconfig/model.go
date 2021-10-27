package currencyconfig

type CurrencyConfig struct {
    Id                      int
    Symbol                  string
    Name                    string
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
    ParentSymbol            string
    ModuleType              string
    Active                  bool
    
    LastUpdated             string
}


