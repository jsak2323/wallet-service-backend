package handlers

import (
	hc "github.com/btcid/wallet-services-backend-go/pkg/domain/healthcheck"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
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
	RpcConfig     RpcConfigResDetail
	HealthCheck   hc.HealthCheck
	IsMaintenance bool
	Error         *errs.Error
}

type GetBlockCountRes struct {
	RpcConfig RpcConfigResDetail
	Blocks    string
	Error     *errs.Error
}

type GetBalanceRes struct {
	RpcConfig RpcConfigResDetail
	Balance   string
	Error     *errs.Error
}

type ListTransactionsRes struct {
	RpcConfig    RpcConfigResDetail
	Transactions []model.Transaction
	Error        *errs.Error
}

type ListWithdrawsRes struct {
    RpcConfig       RpcConfigResDetail
    Withdraws       []model.Withdraw
    Error           string
}

type SendToAddressRes struct {
	RpcConfig RpcConfigResDetail
	TxHash    string
	Error     *errs.Error
}

type GetNewAddressRes struct {
	RpcConfig RpcConfigResDetail
	Address   string
	Error     *errs.Error
}

type AddressTypeRes struct {
	RpcConfig   RpcConfigResDetail
	AddressType string
	Error       *errs.Error
}

type FireblocksSignReq struct {
	Asset       string `json:"asset"`
	DestId      string `json:"destId"`
	DestAddress string `json:"destAddress"`
}

type FireblocksSignRes struct {
	Action          string `json:"action"`
	RejectionReason string `json:"rejectionReason"`
	Error           *errs.Error
}

type StandardRes struct {
	Success bool
	Error   *errs.Error
}
