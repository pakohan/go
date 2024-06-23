package httphelper

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

var PrettyPrint bool

func ServeResponse(w http.ResponseWriter, r *http.Request, err error) {
	// TODO: add known cases
	switch err {
	case nil: // do nothing, will return status 200
	case sql.ErrNoRows:
		http.NotFound(w, r)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func ServeJSON(w http.ResponseWriter, r *http.Request, data interface{}, err error) {
	if err != nil {
		ServeResponse(w, r, err)
		return
	} else if data == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	if PrettyPrint {
		enc.SetIndent("", "  ")
	}

	err = enc.Encode(data)
	if err != nil {
		log.Printf("err in ServeJSON: %v\n", err)
		return
	}
}
