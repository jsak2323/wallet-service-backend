package handlers

import (
    hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
)

const errInternalServer = "Internal server error"

type RpcConfigResDetail struct { 
    RpcConfigId          int
    Symbol               string
    Name                 string
    Host                 string
    Type                 string
    NodeVersion          string
    NodeLastUpdated      string
    IsHealthCheckEnabled bool
}

type GetHealthCheckRes struct { 
    RpcConfig               RpcConfigResDetail
    HealthCheck             hc.HealthCheck
    IsMaintenance           bool
    Error                   string
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

type FireblocksSignReq struct {
	Asset       string `json:"asset"`
	DestAddress string `json:"destAddress"`
}

type FireblocksSignRes struct {
	Action          string `json:"action"`
	RejectionReason string `json:"rejectionReason"`
}

type StandardRes struct {
    Success bool
    Error   string
}


