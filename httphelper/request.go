package httphelper

import (
	"net/http"
	"strconv"
)

type HandleFunc func(*http.Request) (interface{}, error)

func Handle(f HandleFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := f(r)
		ServeJSON(w, r, res, err)
	}
}

func IntParam(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	param, err := strconv.Atoi(r.PathValue(name))
	if err != nil {
		ServeResponse(w, r, err)
		return 0, false
	}

	return param, true
}
