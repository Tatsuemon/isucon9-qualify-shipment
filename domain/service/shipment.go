package service

import "github.com/Tatsuemon/isucon9-qualify-shipment/domain/repository"

type ShipmentService interface {
	Exists(id string) bool
}

type shipmentService struct {
	repository.ShipmentRepository
}

func NewShipmentService(r repository.ShipmentRepository) ShipmentService {
	return &shipmentService{r}
}

func (s *shipmentService) Exists(id string) bool {
	shipment, _ := s.ShipmentRepository.FindByID(id)
	if shipment == nil {
		return true
	}
	return false
}
