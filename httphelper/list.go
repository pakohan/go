package httphelper

import (
	"net/url"
	"strconv"
)

var DefaultPageSize = 10

type List[T interface{}] struct {
	Elements []T `json:"elements"`
	Page     int `json:"page" db:"page"`
	PageSize int `json:"page_size" db:"page_size"`
}

func ListFromQuery[T interface{}](q url.Values) (*List[T], error) {
	res := List[T]{
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
