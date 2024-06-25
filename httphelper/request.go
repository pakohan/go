package httphelper

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// TODO: select parser based on Content-Type
func ParseBody(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil {
		ServeResponse(w, r, err)
		return false
	}

	return true
}

func IntParam(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	param, err := strconv.Atoi(r.PathValue(name))
	if err != nil {
		ServeResponse(w, r, err)
		return 0, false
	}

	return param, true
}

func IntQueryParam(w http.ResponseWriter, r *http.Request, name string) (int, bool) {
	if r.URL.Query().Get(name) == "" {
		return 0, true
	}

	param, err := strconv.Atoi(r.URL.Query().Get(name))
	if err != nil {
		ServeResponse(w, r, err)
		return 0, false
	}

	return param, true
}

func FloatQueryParam(w http.ResponseWriter, r *http.Request, name string) (float64, bool) {
	if r.URL.Query().Get(name) == "" {
		return 0, true
	}

	param, err := strconv.ParseFloat(r.URL.Query().Get(name), 64)
	if err != nil {
		ServeResponse(w, r, err)
		return 0, false
	}

	return param, true
}

func BoolQueryParam(w http.ResponseWriter, r *http.Request, name string) (bool, bool) {
	if r.URL.Query().Get(name) == "" {
		return false, true
	}

	param, err := strconv.ParseBool(r.URL.Query().Get(name))
	if err != nil {
		ServeResponse(w, r, err)
		return false, false
	}

	return param, true
}

func TimeQueryParam(w http.ResponseWriter, r *http.Request, name string) (*time.Time, bool) {
	if r.URL.Query().Get(name) == "" {
		return nil, true
	}

	param, err := time.Parse("2006-01-02T15:04:05", r.URL.Query().Get(name))
	if err != nil {
		ServeResponse(w, r, err)
		return nil, false
	}

	return &param, true
}
