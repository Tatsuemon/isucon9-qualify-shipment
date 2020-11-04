package handler

import (
	"encoding/json"
	"net/http"
)

type AppHandler interface {
	ShipmentHandler
}

// MARK: respose形式
func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJson(w, code, map[string]string{"error": message})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}
