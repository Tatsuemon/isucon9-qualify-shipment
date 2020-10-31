package handler

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Tatsuemon/isucon9-qualify-shipment/domain/entity"
	"github.com/Tatsuemon/isucon9-qualify-shipment/usecase"
	"github.com/skip2/go-qrcode"
)

const IsucariAPIToken = "Bearer 75ugk2m37a750fwir5xr-22l6h4wmue1bwrubzwd0"

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
	toAddress   string `json: to_address`
	toName      string `json: to_name`
	fromAddress string `json: from_address`
	fromName    string `json: from_name`
}

type createRes struct {
	ReserveID   string `json:"reserve_id"`
	ReserveTime int64  `json:"reserve_time`
}

func (s *shipmentHandler) CreateShipment(w http.ResponseWriter, r *http.Request) {
	// TODO(Tatusemon): Middleware
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// TODO(Tatusemon): Middleware
	if r.Header.Get("Authorization") != IsucariAPIToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	req := createReq{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "json decode error")
		return
	}
	defer r.Body.Close()

	ship, err := entity.NewShipment(
		req.toAddress,
		req.toName,
		req.fromAddress,
		req.fromName,
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// TODO(Tatusemon): Middleware
	if r.Header.Get("Authorization") != IsucariAPIToken {
		w.WriteHeader(http.StatusUnauthorized)
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

	// TODO(Tatsuemon): QRの作成 分ける
	schema := "http"
	if r.Header.Get("X-Forwarded-Proto") == "https" {
		schema = "https"
	}
	u := &url.URL{
		Scheme: schema,
		Host:   r.Host,
		Path:   "/accept",
	}
	sha256 := sha256.Sum256([]byte(req.ReserveID))
	q := u.Query()
	q.Set("id", req.ReserveID)
	q.Set("token", fmt.Sprintf("%x", sha256))

	u.RawQuery = q.Encode()

	msg := u.String()

	png, err := qrcode.Encode(msg, qrcode.Low, 256)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h := md5.New()
	h.Write(png)

	_, err := s.ShipmentUseCase.RequestShipment(r.Context(), req.ReserveID, fmt.Sprintf("%x", h.Sum(nil)))
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// TODO(Tatusemon): Middleware
	if r.Header.Get("Authorization") != IsucariAPIToken {
		w.WriteHeader(http.StatusUnauthorized)
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
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// TODO(Tatusemon): Middleware
	if r.Header.Get("Authorization") != IsucariAPIToken {
		w.WriteHeader(http.StatusUnauthorized)
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
