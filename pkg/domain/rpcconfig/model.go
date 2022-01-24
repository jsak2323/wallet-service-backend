package rpcconfig

type RpcConfig struct {
	Id                   int
	Type                 string
	Name                 string `validate:"required"`
	Platform             string `validate:"required"`
	Host                 string `validate:"required"`
	Port                 string
	Path                 string
	User                 string
	Password             string
	Hashkey              string
	NodeVersion          string
	NodeLastUpdated      string
	IsHealthCheckEnabled bool
	AtomFeed             string
	Address              string
	Active               bool
}

type UpdateRpcConfig struct {
	Id                   int `validate:"required"`
	Type                 string
	Name                 string `validate:"required"`
	Platform             string
	Host                 string `validate:"required"`
	Port                 string
	Path                 string `validate:"required"`
	User                 string
	Password             string
	Hashkey              string
	NodeVersion          string
	NodeLastUpdated      string
	IsHealthCheckEnabled bool
	AtomFeed             string
	Address              string
	Active               bool
}
