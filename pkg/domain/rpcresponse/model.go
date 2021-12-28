package rpcresponse

type RpcResponse struct {
	Id              int
	XMLPath         string
	DataTypeXMLTag  string
	TargetFieldName string
	ParseType       string
	JsonFieldsStr   string
	JsonFields      []JsonField
	RpcMethodId     int
}

type JsonField struct {
	SourceFieldName string `json:"source_field_name"`
	TargetFieldName string `json:"target_field_name"`
	DataType        string `json:"data_type"`
}
