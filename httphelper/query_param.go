package httphelper

import (
	"net/http"
	"strconv"
	"time"
)

func IntQueryParam(r *http.Request, name string) (int, error) {
	if r.URL.Query().Get(name) == "" {
		return 0, nil
	}

	param, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil {
		return 0, err
	}

	return param, nil
}

func FloatQueryParam(r *http.Request, name string) (float64, error) {
	if r.URL.Query().Get(name) == "" {
		return 0, nil
	}

	param, err := strconv.ParseFloat(r.URL.Query().Get(name), 64)
	if err != nil {
		return 0, err
	}

	return param, nil
}

func BoolQueryParam(r *http.Request, name string) (bool, error) {
	if r.URL.Query().Get(name) == "" {
		return false, nil
	}

	param, err := strconv.ParseBool(r.URL.Query().Get(name))
	if err != nil {
		return false, err
	}

	return param, nil
}

func TimeQueryParam(r *http.Request, name string) (*time.Time, error) {
	if r.URL.Query().Get(name) == "" {
		return nil, nil
	}

	param, err := time.Parse("2006-01-02T15:04:05", r.URL.Query().Get(name))
	if err != nil {
		return nil, err
	}

	return &param, nil
}
