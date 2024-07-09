package api

import (
	"encoding/json"
	"net/http"
)

func SendJSON[T any](w http.ResponseWriter, val T, status int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(val); err != nil {
		return err
	}
	return nil
}
