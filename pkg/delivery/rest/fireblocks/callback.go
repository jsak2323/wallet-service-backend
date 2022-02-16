package fireblocks

import (
	"encoding/json"
	"net/http"

	"github.com/btcid/wallet-services-backend-go/pkg/http/handlers/fireblocks"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

const TypeBaseAsset = "BASE_ASSET"

const RejectTransaction = "REJECT"
const ApproveTransaction = "APPROVE"
const IgnoreTransaction = "IGNORE"

const InvalidDestAddressReason = "Invalid destination address"

type FireblocksSignReq struct {
	Asset       string `json:"asset" validate:"required"`
	DestId      string `json:"destId" validate:"required"`
	DestAddress string `json:"destAddress" validate:"required"`
}

type FireblocksSignRes struct {
	Action          string `json:"action"`
	RejectionReason string `json:"rejectionReason"`
	Error           *errs.Error
}

func (re *Rest) CallbackHandler(w http.ResponseWriter, req *http.Request) {
	var (
		SignReq fireblocks.FireblocksSignReq
		RES     fireblocks.FireblocksSignRes
		err     error
		ctx     = req.Context()
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError

			RES.Action = RejectTransaction
			RES.RejectionReason = RES.Error.Error()

			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		} else {
			logger.InfoLog(" -- fireblocks.CallbackHandler Success!", req)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" -- fireblocks.CallbackHandler, Requesting ...", req)

	RES.Action = ApproveTransaction

	if err = json.NewDecoder(req.Body).Decode(&SignReq); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.ErrorUnmarshalBodyRequest)
		return
	}

	service := re.svc.Fireblocks
	if RES, err = service.ValidateHotDestAddress(ctx, SignReq); err != nil {
		RES.Error = errs.AddTrace(err)
		return
	}

}
