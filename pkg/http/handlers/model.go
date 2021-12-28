package handlers

import (
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

const errInternalServer = "Internal server error"
const errAssetNotFound = "Asset not found"

type RpcConfigResDetail struct { 
    RpcConfigId          int
    Symbol               string
    TokenType            string
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
    Transactions    []model.Transaction
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
	Type        string `json:"type"`
    DestId      string `json:"destId"`
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


