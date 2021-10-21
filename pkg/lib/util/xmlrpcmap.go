package util

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	mxj "github.com/clbanning/mxj/v2"
	"github.com/divan/gorilla-xmlrpc/xml"

	rr "github.com/btcid/wallet-services-backend-go/pkg/domain/rpcresponse"
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

func (xrm *XmlRpcMap) XmlRpcMapCall(method string, args *XmlRpcMapReq, resFieldMap map[string]rr.RpcResponse, reply model.RpcRes) error {
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

	mapValues, err := ParseResponse(res.Body, resFieldMap)
	if err != nil {
		fmt.Println(" - xmlrpcmap.ParseResponse err: " + err.Error())
		return err
	}

	if err = reply.SetFromMapValues(mapValues); err != nil {
		fmt.Println(" - reply.SetFromMapValues err: " + err.Error())
		return err
	}

	return nil
}

func ParseResponse(resBody io.ReadCloser, rpcResMap map[string]rr.RpcResponse) (resValues map[string]interface{}, err error) {
	mv, err := mxj.NewMapXmlReader(resBody)
	if err != nil {
		return map[string]interface{}{}, err
	}

	resValues = make(map[string]interface{})

	for _, rpcRes := range rpcResMap {
		pathArr := strings.Split(rpcRes.XMLPath, ".")

		for _, tag := range pathArr {
			var ok bool

			mv, ok = mv[tag].(map[string]interface{})

			if !ok && rpcRes.FieldName != rr.FieldNameError {
				return map[string]interface{}{}, errors.New("mismatched rpc response config")
			}
		}

		resValues[rpcRes.FieldName] = mv[rpcRes.DataTypeTag]
	}

	return resValues, nil
}
