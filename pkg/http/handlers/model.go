package handlers

type GetBlockCountRes struct{
    RpcConfigId         int
    Symbol              string
    Name                string
    Host                string
    Type                string
    NodeVersion         string
    NodeLastUpdated     string
    Blocks              string
}

type GetBalanceRes struct {
    RpcConfigId         int
    Symbol              string
    Name                string
    Host                string
    Type                string
    Balance             string
}