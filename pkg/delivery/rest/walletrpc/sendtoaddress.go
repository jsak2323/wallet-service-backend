package walletrpc

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type SendToAddressRequest struct {
	Symbol    string `json:"symbol"`
	TokenType string `json:"token_type"`
	Amount    string `json:"amount"`
	Address   string `json:"address"`
	Memo      string `json:"memo"`
}

func (re *Rest) SendToAddressHandler(w http.ResponseWriter, req *http.Request) {
	// define response object
	RES := handlers.SendToAddressRes{}
	ctx := req.Context()

	// define response handler
	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	// define request params
	sendToAddressRequest := handlers.SendToAddressRequest{}
	err := DecodeAndLogPostRequest(req, &sendToAddressRequest)
	if err != nil {
		logger.ErrorLog(" - SendToAddressHandler util.DecodeAndLogPostRequest(req, &sendToAddressRequest) err: "+err.Error(), ctx)
		return
	}

	service := re.svc.WalletRpc
	res, err := service.SendToAddress(ctx, sendToAddressRequest)
	if err != nil {
		RES.Error = errs.AddTrace(errors.New(res.Error.Error()))
	} else {
		resJson, _ := json.Marshal(RES)
		logger.InfoLog(" - SendToAddressHandler Success. Symbol: "+strings.ToUpper(sendToAddressRequest.Symbol)+", Res: "+string(resJson), req)
	}

}

func DecodeAndLogPostRequest(req *http.Request, output interface{}) error {
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return errs.AddTrace(err)
	}

	logger.InfoLog("POST Request Body : "+string(reqBody), req)

	err = json.Unmarshal(reqBody, output)
	if err != nil {
		return errs.AddTrace(err)
	}

	return nil
}
