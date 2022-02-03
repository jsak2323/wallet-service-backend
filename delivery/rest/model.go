package rest

import (
	"context"
	"encoding/json"
	"net/http"

	domain "github.com/btcid/wallet-services-backend-go/pkg/domain/permission"
	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

type StandardRes struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}

type CreateReq struct {
	Name string `json:"name"`
}

type CreateRes struct {
	Id      int         `json:"id"`
	Message string      `json:"message"`
	Error   *errs.Error `json:"error"`
}

type ListRes struct {
	Permissions []domain.Permission `json:"permissions"`
	Error       *errs.Error         `json:"error"`
}

func HandleResponse(ctx context.Context, w http.ResponseWriter, resp interface{}, err *errs.Error) {
	resStatus := http.StatusOK
	if err != nil {
		resStatus = http.StatusInternalServerError
		logger.ErrorLog(errs.Logged(err), ctx)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resStatus)
	json.NewEncoder(w).Encode(resp)
}
