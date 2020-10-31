package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// TODO(Tatsuemon): enum的な何か
// type ShipmentStatus struct {
// 	string
// }

// shipmentStatus := [...]string {"initial", "wait_pickup", "shipping", "done", "cancel"}

// func NewShipmentStatus() (*ShipmentStatus, error) {
// 	return &ShipmentStatus{shipmentStatus[0]}
// }

type Shipment struct {
	ID              string    `json:"id" db:"id"`
	ToAddress       string    `json:"to_address" db:"to_address"`
	ToName          string    `json:"to_name" db:"to_name"`
	FromAddress     string    `json:"from_address" db:"from_address"`
	FromName        string    `json:"from_name" db:"from_name"`
	Status          string    `json:"status" db:"status"`
	ReserveDateTime time.Time `json:"reseve_date_time" db:"reserve_date_time"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

func NewShipment(toAddress string, toName string, fromAddress string, fromName string) (*Shipment, error) {
	err := validateShipment(toAddress, toName, fromAddress, fromName)
	if err != nil {
		return nil, err
	}
	return &Shipment{
		ID:              uuid.New().String(),
		ToAddress:       toAddress,
		ToName:          toName,
		FromAddress:     fromAddress,
		FromName:        fromName,
		Status:          "initial",
		ReserveDateTime: time.Now().AddDate(0, 0, 1), // 1日後
		CreatedAt:       time.Now(),
	}, nil
}

func validateShipment(toAddress string, toName string, fromAddress string, fromName string) error {
	if toAddress == "" {
		return fmt.Errorf("toAddress is required.")
	}
	if toName == "" {
		return fmt.Errorf("toName is required.")
	}
	if fromAddress == "" {
		return fmt.Errorf("fromAddress is required.")
	}
	if fromName == "" {
		return fmt.Errorf("fromName is required.")
	}
	return nil
}

// 東京のタイムゾーンを取得
//
// @return *time.Location
func Location() *time.Location {
	return time.FixedZone("Asia/Tokyo", 9*60*60)
}
