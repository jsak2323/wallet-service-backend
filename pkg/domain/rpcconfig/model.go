package rpcconfig

type RpcConfig struct {
    Id                      int
    CurrencyId              int
    Type                    string
    Name                    string
    Platform                string
    Host                    string
    Port                    string
    Path                    string
    User                    string
    Password                string
    Hashkey                 string
    NodeVersion             string
    NodeLastUpdated         string
    IsBackupable            bool
    IsHealthCheckEnabled    bool
    AtomFeed                string
}


