package handlers

import (
    hc "github.com/btcid/wallet-services-backend/pkg/domain/healthcheck"
)

type RpcConfigResDetail struct { 
    RpcConfigId         int
    Symbol              string
    Name                string
    Host                string
    Type                string
    NodeVersion         string
    NodeLastUpdated     string
}

type GetHealthCheckRes struct { 
    RpcConfig    RpcConfigResDetail
    HealthCheck  hc.HealthCheck
    Error        string
}

type GetBlockCountRes struct { 
    RpcConfig   RpcConfigResDetail
    Blocks      string
    Error       string
}

type GetBalanceRes struct { 
    RpcConfig   RpcConfigResDetail
    Balance     string
    Error       string
}

type ListTransactionsRes struct {
    RpcConfig       RpcConfigResDetail
    Transactions    string
    Error           string
}

type SendToAddressRes struct {
    RpcConfig   RpcConfigResDetail
    TxHash      string
    Error       string
}

type GetNewAddressRes struct {
    RpcConfig   RpcConfigResDetail
    Address     string
    Error       string
}

type AddressTypeRes struct {
    RpcConfig   RpcConfigResDetail
    AddressType string
    Error       string
}


