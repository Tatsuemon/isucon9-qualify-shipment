package router

import (
	"github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/web/handler"
	"github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/web/middleware"
	"github.com/gorilla/mux"
)

func NewMuxRouter(h handler.AppHandler) *mux.Router {
	mux := mux.NewRouter()
	mux.Use(middleware.AuthrizationBearer)
	mux.HandleFunc("/create", h.CreateShipment).Methods("POST")
	mux.HandleFunc("/request", h.RequestShipment).Methods("POST")
	mux.HandleFunc("/done", h.DoneShipment).Methods("POST")
	mux.HandleFunc("/accept", h.AcceptShipment).Methods("GET")
	mux.HandleFunc("/status", h.StatusShipment).Methods("GET")
	return mux
}
