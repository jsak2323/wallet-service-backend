package rpcresponse

type CreateRpcResponse struct {
	XMLPath         string `validate:"required"`
	DataTypeXMLTag  string `validate:"required"`
	TargetFieldName string `validate:"required"`
	ParseType       string `validate:"required"`
	JsonFieldsStr   string
	JsonFields      []JsonField
	RpcMethodId     int `validate:"required"`
}

type RpcResponse struct {
	Id int `validate:"required"`
	CreateRpcResponse
}

type JsonField struct {
	SourceFieldName string `json:"source_field_name"`
	TargetFieldName string `json:"target_field_name"`
	DataType        string `json:"data_type"`
}
