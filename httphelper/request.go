package httphelper

import (
	"encoding/json"
	"net/http"
	"strconv"
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
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ServeResponse(w, r, err)
		return 0, false
	}

	return id, true
}
