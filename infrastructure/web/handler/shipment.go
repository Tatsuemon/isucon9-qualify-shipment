package handler

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Tatsuemon/isucon9-qualify-shipment/domain/entity"
	"github.com/Tatsuemon/isucon9-qualify-shipment/usecase"
)

var isucariAPIToken = os.Getenv("AUTH_BEARER")

type ShipmentHandler interface {
	CreateShipment(w http.ResponseWriter, r *http.Request)
	RequestShipment(w http.ResponseWriter, r *http.Request)
	AcceptShipment(w http.ResponseWriter, r *http.Request)
	StatusShipment(w http.ResponseWriter, r *http.Request)
	DoneShipment(w http.ResponseWriter, r *http.Request)
}

type shipmentHandler struct {
	usecase.ShipmentUseCase
}

func NewShipmentHandler(s usecase.ShipmentUseCase) ShipmentHandler {
	return &shipmentHandler{s}
}

// [POST] /create
type createReq struct {
	ToAddress   string `json:"to_address"`
	ToName      string `json:"to_name"`
	FromAddress string `json:"from_address"`
	FromName    string `json:"from_name"`
}

type createRes struct {
	ReserveID   string `json:"reserve_id"`
	ReserveTime int64  `json:"reserve_time"`
}

func (s *shipmentHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	// TODO(Tatusemon): Middleware
	if r.Header.Get("Authorization") != isucariAPIToken {
		respondWithError(w, http.StatusUnauthorized, "authorization error")
		return
	}

	req := createReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "json decode error")
		return
	}
	defer r.Body.Close()

	ship, err := entity.NewShipment(
		req.ToAddress,
		req.ToName,
		req.FromAddress,
		req.FromName,
	)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "required parameter was not passed")
		return
	}

	now := time.Now()
	ship.ReserveDateTime = now
	ship, err = s.ShipmentUseCase.CreateShipment(r.Context(), ship)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	res := createRes{}
	res.ReserveID = ship.ID
	res.ReserveTime = ship.ReserveDateTime.Unix()

	respondWithJson(w, http.StatusOK, res)
	return
}

// [POST] /request
type requestReq struct {
	ReserveID string `json:"reserve_id"`
}

func (s *shipmentHandler) RequestShipment(w http.ResponseWriter, r *http.Request) {
	// TODO(Tatusemon): Middleware
	if r.Header.Get("Authorization") != isucariAPIToken {
		respondWithError(w, http.StatusUnauthorized, "authorization error")
		return
	}

	req := requestReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "json decode error")
		return
	}

	if req.ReserveID == "" {
		respondWithError(w, http.StatusBadRequest, "required parameter was not passed")
		return
	}

	schema := r.Header.Get("X-Forwarded-Proto")
	host := r.Host
	png, err := s.ShipmentUseCase.CreateAcceptQr(schema, host, req.ReserveID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h := md5.New()
	h.Write(png)

	_, err = s.ShipmentUseCase.RequestShipment(r.Context(), req.ReserveID, fmt.Sprintf("%x", h.Sum(nil)))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "empty")
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(png)
}

// [GET] /accept
func (s *shipmentHandler) AcceptShipment(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	id := query.Get("id")
	token := query.Get("token")

	if ok := s.ShipmentUseCase.CheckAcceptToken(id, token); !ok {
		respondWithError(w, http.StatusBadRequest, "wrong parameters")
		return
	}

	_, err := s.ShipmentUseCase.SetShippingStatus(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "empty")
		return
	}
	response := map[string]interface{}{"accept": "ok"}
	respondWithJson(w, http.StatusOK, response)
	return
}

// [GET] /status
type shipmentStatusReq struct {
	ReserveID string `json:"reserve_id"`
}

type shipmentStatusRes struct {
	Status      string `json:"status"`
	ReserveTime int64  `json:"reserve_time"`
}

func (s *shipmentHandler) StatusShipment(w http.ResponseWriter, r *http.Request) {
	// TODO(Tatusemon): Middleware
	if r.Header.Get("Authorization") != isucariAPIToken {
		respondWithError(w, http.StatusUnauthorized, "authorization error")
		return
	}

	req := shipmentStatusReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "json decode error")
		return
	}

	if req.ReserveID == "" {
		respondWithError(w, http.StatusBadRequest, "required parameter was not passed")
		return
	}
	status, reserveTime, err := s.ShipmentUseCase.GetStatus(req.ReserveID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	res := shipmentStatusRes{}
	res.Status = status
	res.ReserveTime = reserveTime
	respondWithJson(w, http.StatusOK, res)
	return
}

// [POST] /done
type doneShipmentReq struct {
	ReserveID string `json:"reserve_id"`
}

func (s *shipmentHandler) DoneShipment(w http.ResponseWriter, r *http.Request) {
	// TODO(Tatusemon): Middleware
	if r.Header.Get("Authorization") != isucariAPIToken {
		respondWithError(w, http.StatusUnauthorized, "authorization error")
		return
	}

	req := doneShipmentReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "json decode error")
		return
	}

	if req.ReserveID == "" {
		respondWithError(w, http.StatusBadRequest, "required parameter was not passed")
		return
	}

	_, err := s.ShipmentUseCase.DoneShipment(r.Context(), req.ReserveID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "empty")
		return
	}
	response := map[string]interface{}{"accept": "ok"}
	respondWithJson(w, http.StatusOK, response)
	return
}
