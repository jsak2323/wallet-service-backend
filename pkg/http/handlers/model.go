package handlers

type RpcConfigResDetail struct { 
    RpcConfigId         int
    Symbol              string
    Name                string
    Host                string
    Type                string
    NodeVersion         string
    NodeLastUpdated     string
}

type GetBlockCountRes struct { 
    RpcConfig   RpcConfigResDetail
    Blocks      string
}

type GetBalanceRes struct { 
    RpcConfig   RpcConfigResDetail
    Balance     string
}

type ListTransactionsRes struct {
    RpcConfig       RpcConfigResDetail
    Transactions    string
}

type SendToAddressRes struct {
    RpcConfig   RpcConfigResDetail
    TxHash      string
}