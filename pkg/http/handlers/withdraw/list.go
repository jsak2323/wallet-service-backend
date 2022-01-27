package withdraw

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	errs "github.com/btcid/wallet-services-backend-go/pkg/lib/error"
	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *WithdrawService) ListHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES ListRes
		err error
		ctx = req.Context()

		page, limit int
		filters     []map[string]interface{}
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		}

		resStatus = http.StatusOK
		if RES.Error != nil {
			resStatus = http.StatusInternalServerError
			logger.ErrorLog(errs.Logged(RES.Error), ctx)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - withdraw.ListHandler, Requesting ...", req)

	if filters, err = parseFilters(req); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedParseFilter)
		return
	}

	if page, limit, err = parsePagination(req); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedParsePagination)
		return
	}

	if RES.Withdraws, err = s.wRepo.Get(page, limit, filters); err != nil {
		RES.Error = errs.AssignErr(errs.AddTrace(err), errs.FailedGetListWithdraws)
		return
	}
}

func parseFilters(req *http.Request) (filters []map[string]interface{}, err error) {
	var currencyId int

	filters = make([]map[string]interface{}, 0)

	if req.URL.Query().Get("currency_id") != "" {
		if currencyId, err = strconv.Atoi(req.URL.Query().Get("currency_id")); err != nil {
			err = errs.AddTrace(errors.New("error parsing currency_id value: " + req.URL.Query().Get("currency_id") + " Error: " + err.Error()))
			return []map[string]interface{}{}, err
		}

		filters = append(filters, map[string]interface{}{"key": "currency_id", "value": currencyId})
	}

	if req.URL.Query().Get("address_to") != "" {
		filters = append(filters, map[string]interface{}{"key": "address_to", "value": req.URL.Query().Get("address_to")})
	}

	if req.URL.Query().Get("tx") != "" {
		filters = append(filters, map[string]interface{}{"key": "tx", "value": req.URL.Query().Get("tx")})
	}

	if req.URL.Query().Get("success_date") != "" {
		filters = append(filters, map[string]interface{}{"key": "success_time", "value": req.URL.Query().Get("success_date"), "format": "date"})
	}

	return filters, nil
}

func parsePagination(req *http.Request) (page, limit int, err error) {
	if page, err = strconv.Atoi(req.URL.Query().Get("page")); err != nil && req.URL.Query().Get("page") != "" {
		err = errs.AddTrace(errors.New("error parsing page value: " + req.URL.Query().Get("page") + " Error: " + err.Error()))
		return 0, 0, err
	}

	if limit, err = strconv.Atoi(req.URL.Query().Get("limit")); err != nil && req.URL.Query().Get("limit") != "" {
		err = errs.AddTrace(errors.New("error parsing limit value: " + req.URL.Query().Get("limit") + " Error: " + err.Error()))
		return 0, 0, err
	}

	return page, limit, nil
}
