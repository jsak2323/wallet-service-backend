package rpcresponse

import (
	"errors"
	"fmt"
	"encoding/json"
)

func(r RpcResponse) ParseField(source, target interface{}) (ok bool) {
	switch r.ParseType {
	case "string":
		target, ok = source.(string)
	case "int":
		target, ok = source.(int)
	default:
		return false
	}
	
	return ok
}

func(r RpcResponse) ParseArrayOfJson(source interface{}) (result []map[string]interface{}, err error) {
	var ok bool
	var tempString string
	var arrayOfJson []map[string]interface{}

	if tempString, ok = source.(string); !ok {
		return nil, errors.New("mismatched rpc response parse type and data type")
	}

	if err = json.Unmarshal([]byte(tempString), &arrayOfJson); err != nil {
		return nil, err
	}

	for _, jsonData := range arrayOfJson {
		targetJsonMap := map[string]interface{}{}
		for _, jsonField := range r.JsonFields {
			switch jsonField.DataType {
				case "string": {
					targetJsonMap[jsonField.TargetFieldName] = jsonData[jsonField.SourceFieldName].(string)
				}
				case "int64": {
					targetJsonMap[jsonField.TargetFieldName] = jsonData[jsonField.SourceFieldName].(int64)
				}
				default: {
					return nil, fmt.Errorf("Invalid target json field data_type: %s rpc_response_id: %d", jsonField.DataType, r.Id)
				}
			}
		}
		result = append(result, targetJsonMap)
	}

	return result, nil
}