package usecase

import (
	"context"
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/Tatsuemon/isucon9-qualify-shipment/domain/repository"

	"github.com/Tatsuemon/isucon9-qualify-shipment/domain/entity"
	"github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/datastore"
)

type ShipmentUseCase interface {
	GetStatus(id string) (string, int64, error)                                              // GET /status
	SetShippingStatus(ctx context.Context, id string) (*entity.Shipment, error)              // GET /accept
	CreateShipment(ctx context.Context, shipment *entity.Shipment) (*entity.Shipment, error) // POST /create
	RequestShipment(ctx context.Context, id string, str string) (*entity.Shipment, error)    // POST /request
	DoneShipment(ctx context.Context, id string) (*entity.Shipment, error)                   // POST /done
	CheckAcceptToken(id string, token string) bool
}

type shipmentUseCase struct {
	repository.ShipmentRepository
	transaction datastore.Transaction
}

func NewShipmentUseCase(r repository.ShipmentRepository, t datastore.Transaction) ShipmentUseCase {
	return &shipmentUseCase{r, t}
}

func (s *shipmentUseCase) GetStatus(id string) (string, int64, error) {
	shipment, err := s.ShipmentRepository.FindByID(id)
	if err != nil {
		return "", 0, err
	}
	return shipment.Status, shipment.ReserveDateTime.Unix(), nil
}

func (s *shipmentUseCase) SetShippingStatus(ctx context.Context, id string) (*entity.Shipment, error) {
	v, err := s.transaction.DoInTx(ctx, func(context.Context) (interface{}, error) {
		shipment, err := s.ShipmentRepository.FindByID(id)
		if err != nil {
			return nil, err
		}
		shipment.Status = entity.SHIPPING
		return s.ShipmentRepository.Update(ctx, shipment)
	})
	if err != nil {
		return nil, err
	}
	return v.(*entity.Shipment), nil
}

func (s *shipmentUseCase) CreateShipment(ctx context.Context, shipment *entity.Shipment) (*entity.Shipment, error) {
	v, err := s.transaction.DoInTx(ctx, func(context.Context) (interface{}, error) {
		return s.ShipmentRepository.Store(ctx, shipment)
	})
	if err != nil {
		return nil, err
	}
	return v.(*entity.Shipment), nil
}

func (s *shipmentUseCase) RequestShipment(ctx context.Context, id string, str string) (*entity.Shipment, error) {
	v, err := s.transaction.DoInTx(ctx, func(context.Context) (interface{}, error) {
		shipment, err := s.ShipmentRepository.FindByID(id)
		if err != nil {
			return nil, err
		}
		shipment.Status = entity.WAIT_PICKUP
		shipment.QRMD5 = str
		return s.ShipmentRepository.Update(ctx, shipment)
	})
	if err != nil {
		return nil, err
	}
	return v.(*entity.Shipment), nil
}

func (s *shipmentUseCase) DoneShipment(ctx context.Context, id string) (*entity.Shipment, error) {
	v, err := s.transaction.DoInTx(ctx, func(context.Context) (interface{}, error) {
		shipment, err := s.ShipmentRepository.FindByID(id)
		if err != nil {
			return nil, err
		}
		shipment.Status = entity.DONE
		shipment.DoneDateTime = time.Now()
		return s.ShipmentRepository.Update(ctx, shipment)
	})
	if err != nil {
		return nil, err
	}
	return v.(*entity.Shipment), nil
}

func (s *shipmentUseCase) CheckAcceptToken(id string, token string) bool {
	sha256 := sha256.Sum256([]byte(id))
	return token == fmt.Sprintf("%x", sha256)
}
