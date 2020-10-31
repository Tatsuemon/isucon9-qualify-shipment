package repository

import (
	"context"

	"github.com/Tatsuemon/isucon9-qualify-shipment/domain/entity"
)

type ShipmentRepository interface {
	FindByID(id string) (*entity.Shipment, error)
	Store(ctx context.Context, shipment *entity.Shipment) (*entity.Shipment, error)
	Update(ctx context.Context, shipment *entity.Shipment) (*entity.Shipment, error)
}
