package api

import (
	"encoding/json"
	"net/http"
)

func WriteErr(w http.ResponseWriter, errCode int, data ...interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(errCode)
	json.NewEncoder(w).Encode(data)
}
