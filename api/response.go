package api

import (
	"encoding/json"
	"net/http"
)

func Write(w http.ResponseWriter, data interface{}) {

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(200)

	json.NewEncoder(w).Encode(data)
}
