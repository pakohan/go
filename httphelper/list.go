package httphelper

import (
	"net/url"
	"strconv"
)

var DefaultPageSize = 10

type List[T, F any] struct {
	Elements []T `json:"elements"`
	Filter   F   `json:"filter"`
	Page     int `json:"page" db:"page"`
	PageSize int `json:"page_size" db:"page_size"`
}

func ListFromQuery[T, F any](q url.Values) (*List[T, F], error) {
	res := List[T, F]{
		Page:     1,
		PageSize: DefaultPageSize,
	}

	ps := q.Get("page")
	if ps != "" {
		var err error
		res.Page, err = strconv.Atoi(ps)
		if err != nil {
			return nil, err
		}
	}

	pss := q.Get("page_size")
	if pss != "" {
		var err error
		res.PageSize, err = strconv.Atoi(pss)
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}
