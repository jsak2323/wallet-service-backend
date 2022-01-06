package mysql

import (
	"github.com/mitchellh/mapstructure"
)

type Filter struct {
	Key      string
	Value    interface{}
	Format   string
	Operator string
}

func parseFilters(filters []map[string]interface{}, query *string, params *[]interface{}) error {
	for i, mapFilter := range filters {
		filter := Filter{}

		if err := mapstructure.Decode(mapFilter, &filter); err != nil {
			return err
		}

		if i == 0 {
			*query = *query + " WHERE "
		} else {
			*query = *query + " AND "
		}

		if filter.Operator == "" {
			filter.Operator = "="
		}

		if filter.Format == "date" {
			*query = *query + "DATE(" + filter.Key + ") " + filter.Operator + " DATE_FORMAT(?, '%Y-%m-%d')"
		} else {
			*query = *query + filter.Key + " " + filter.Operator + " ? "
		}

		*params = append(*params, filter.Value)
	}

	return nil
}
