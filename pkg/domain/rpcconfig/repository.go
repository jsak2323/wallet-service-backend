package rpcconfig

type RpcConfigRepository interface {
    GetByCurrencyId(currency_id int) ([]RpcConfig, error)
}