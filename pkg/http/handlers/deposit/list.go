package deposit

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	logger "github.com/btcid/wallet-services-backend-go/pkg/logging"
)

func (s *DepositService) ListHandler(w http.ResponseWriter, req *http.Request) {
	var (
		RES ListRes
		err error

		page, limit int
		filters     []map[string]interface{}
	)

	handleResponse := func() {
		resStatus := http.StatusOK
		if err != nil {
			resStatus = http.StatusInternalServerError
		}
		w.WriteHeader(resStatus)
		json.NewEncoder(w).Encode(RES)
	}
	defer handleResponse()

	logger.InfoLog(" - deposit.ListHandler, Requesting ...", req)

	if filters, err = parseFilters(req); err != nil {
		logger.ErrorLog(" -- deposit.ListHandler parseFilters Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if page, limit, err = parsePagination(req); err != nil {
		logger.ErrorLog(" -- deposit.ListHandler parsePagination Error: " + err.Error())
		RES.Error = err.Error()
		return
	}

	if RES.Deposits, err = s.dRepo.Get(page, limit, filters); err != nil {
		logger.ErrorLog(" -- deposit.ListHandler dRepo.Get Error: " + err.Error())
		RES.Error = err.Error()
		return
	}
}

func parseFilters(req *http.Request) (filters []map[string]interface{}, err error) {
	var currencyId int

	filters = make([]map[string]interface{}, 0)

	if req.URL.Query().Get("currency_id") != "" {
		if currencyId, err = strconv.Atoi(req.URL.Query().Get("currency_id")); err != nil {
			err = fmt.Errorf("error parsing currency_id value: " + req.URL.Query().Get("currency_id") + " Error: " + err.Error())
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
		err = fmt.Errorf("error parsing page value: " + req.URL.Query().Get("page") + " Error: " + err.Error())
		return 0, 0, err
	}

	if limit, err = strconv.Atoi(req.URL.Query().Get("limit")); err != nil && req.URL.Query().Get("limit") != "" {
		err = fmt.Errorf("error parsing limit value: " + req.URL.Query().Get("limit") + " Error: " + err.Error())
		return 0, 0, err
	}

	return page, limit, nil
}
