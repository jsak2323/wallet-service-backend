package handlers

type GetBlockCountRes struct{
    Symbol              string
    Host                string
    Type                string
    NodeVersion         string
    NodeLastUpdated     string
    Blocks              string
}