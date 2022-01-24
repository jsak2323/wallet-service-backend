package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	mxj "github.com/clbanning/mxj/v2"
	"github.com/divan/gorilla-xmlrpc/xml"

	rc "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcconfig"
	rm "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcmethod"
	rrq "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcrequest"
	rrs "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
	"github.com/btcid/wallet-services-backend-go/pkg/modules/model"
)

// Arbitrary num of fields
type XmlRpcMapReq struct {
	Arg1 string
	Arg2 string
	Arg3 string
	Arg4 string
	Arg5 string
	Arg6 string
	Arg7 string
	Arg8 string
}

type XmlRpcMap struct {
	Host string
	Port string
	Path string
}

func NewXmlRpcMapClient(host string, port string, path string) *XmlRpcMap {
	return &XmlRpcMap{
		host, port, path,
	}
}

func (xrm *XmlRpcMap) XmlRpcMapCall(method string, args *XmlRpcMapReq, resFieldMap map[string]rrs.RpcResponse, reply model.RpcRes) error {
	buf, err := xml.EncodeClientRequest(method, args)
	if err != nil {
		fmt.Println(" - xml.EncodeClientRequest(method, args) err: " + err.Error())
		return err
	}

	url := "http://" + xrm.Host + ":" + xrm.Port + xrm.Path
	httpClient := &http.Client{
		Timeout: 120 * time.Second,
	}
	res, err := httpClient.Post(url, "text/xml", bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(" - httpClient.Post xml err: " + err.Error())
		return err
	}
	defer res.Body.Close()

	mapValues, err := DecodeResponseToMap(res.Body, resFieldMap)
	if err != nil {
		fmt.Println(" - xmlrpcmap.ParseResponse err: " + err.Error())
		return err
	}

	if err = reply.SetFromMapValues(mapValues, resFieldMap); err != nil {
		fmt.Println(" - reply.SetFromMapValues err: " + err.Error())
		return err
	}

	return nil
}

func RpcRequestValue(rpcRequest rrq.RpcRequest, runtimeParams map[string]string) (string, error) {
	if rpcRequest.Source == rrq.SourceRuntime {
		runtimeParam, ok := runtimeParams[rpcRequest.RuntimeVarName]
		if !ok {
			return "", errors.New("runtime param not passed: " + rpcRequest.RuntimeVarName)
		}

		return runtimeParam, nil
	}

	if rpcRequest.Source == rrq.SourceConfig {
		return rpcRequest.Value, nil
	}

	return "", errors.New("invalid rpc request source: " + rpcRequest.Source)
}

func RpcRequestJsonField(jsonRpcRequests []rrq.RpcRequest, runtimeParams map[string]string) (string, error) {
	var err error

	fieldMap := make(map[string]interface{})
	for _, field := range jsonRpcRequests {
		if fieldMap[field.ArgName], err = RpcRequestValue(field, runtimeParams); err != nil {
			return "", err
		}
	}

	jsonEncoded, err := json.Marshal(fieldMap)
	if err != nil {
		return "", err
	}

	return string(jsonEncoded), nil
}

func GetRpcRequestArgs(rpcConfig rc.RpcConfig, rpcMethod rm.RpcMethod, rpcRequests []rrq.RpcRequest, runtimeParams map[string]string) (args []string, err error) {
	args = make([]string, rpcMethod.NumOfArgs)

	hashkey, nonce := GenerateHashkey(rpcConfig.Password, rpcConfig.Hashkey)

	if len(args) >= 1 {
		args[0] = rpcConfig.User
	}
	if len(args) >= 2 {
		args[1] = hashkey
	}
	if len(args) >= 5 {
		args[5] = nonce
	}

	for _, rpcRequest := range rpcRequests {
		switch rpcRequest.Type {
		case rrq.TypeJsonRoot:
			args[rpcRequest.ArgOrder], err = RpcRequestJsonField(rrq.JsonFieldTypeRpcRequests(rpcRequests), runtimeParams)
			if err != nil {
				return []string{}, err
			}
		case rrq.TypeValueRoot:
			args[rpcRequest.ArgOrder], err = RpcRequestValue(rpcRequest, runtimeParams)
			if err != nil {
				return []string{}, err
			}
		default:
			continue
		}
	}

	return args, nil
}

func GenerateRpcMapRequest(args []string) (req XmlRpcMapReq) {
	for i := 0; i < reflect.ValueOf(&req).Elem().NumField(); i++ {
		if i >= len(args) {
			break
		}

		reqField := reflect.ValueOf(&req).Elem().Field(i)
		reqField.SetString(args[i])
	}

	return req
}

func DecodeResponseToMap(resBody io.ReadCloser, rpcResMap map[string]rrs.RpcResponse) (resValues map[string]interface{}, err error) {
	xmlResMap, err := mxj.NewMapXmlReader(resBody)
	if err != nil {
		return map[string]interface{}{}, err
	}

	resValues = make(map[string]interface{})

	for _, rpcRes := range rpcResMap {
		// skip error first
		if rpcRes.TargetFieldName == rrs.FieldNameError {
			continue
		}

		if resValues[rpcRes.TargetFieldName], err = getValueFromResMap(rpcRes, xmlResMap); err != nil {
			return map[string]interface{}{}, err
		}
	}

	// configured path not matching xmlResMap: could be misconfiguration or response returned error
	if err != nil {
		// find error xml path in xmlResMap
		errRpcRes, ok := rpcResMap[rrs.FieldNameError]
		if !ok {
			return map[string]interface{}{}, fmt.Errorf("rpc_method have no error field rpc_response %+v", rpcResMap)
		}

		if resValues[errRpcRes.TargetFieldName], err = getValueFromResMap(errRpcRes, xmlResMap); err != nil {
			return map[string]interface{}{}, err
		}
	}

	return resValues, nil
}

func getValueFromResMap(rpcRes rrs.RpcResponse, xmlResMap map[string]interface{}) (value interface{}, err error) {
	pathArr := strings.Split(rpcRes.XMLPath, ".")

	for _, tag := range pathArr {
		var ok bool

		xmlResMap, ok = xmlResMap[tag].(map[string]interface{})
		if !ok {
			err = errors.New("mismatched rpc response config: " + rpcRes.TargetFieldName)
			break
		}
	}

	return xmlResMap[rpcRes.DataTypeXMLTag], nil
}
