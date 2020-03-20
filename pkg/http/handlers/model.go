package handlers

type GetBlockCountRes struct{
    RpcConfigId 		int
    Symbol              string
    Host                string
    Type                string
    NodeVersion         string
    NodeLastUpdated     string
    Blocks              string
}