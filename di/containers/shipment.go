package di

import (
	"log"

	"github.com/Tatsuemon/isucon9-qualify-shipment/domain/repository"
	"github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/datastore"
	"github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/web/handler"
	"github.com/Tatsuemon/isucon9-qualify-shipment/usecase"
	"github.com/jmoiron/sqlx"
)

type ShipmentContainer struct {
	Repository repository.ShipmentRepository
	Usecase    usecase.ShipmentUseCase
	Handler    handler.ShipmentHandler
}

func NewShipmentContainer(env string, db *sqlx.DB) *ShipmentContainer {
	switch env {
	case "development":
		log.Print("development")
		shipmentTransaction := datastore.NewTransaction(db)
		shipmentRepository := datastore.NewShipmentRepository(db)
		shipmentUsecase := usecase.NewShipmentUseCase(shipmentRepository, shipmentTransaction)
		shipmentHandler := handler.NewShipmentHandler(shipmentUsecase)
		return &ShipmentContainer{
			Repository: shipmentRepository,
			Usecase:    shipmentUsecase,
			Handler:    shipmentHandler,
		}
	case "test":
		log.Print("test")
		shipmentTransaction := datastore.NewTransaction(db)
		shipmentRepository := datastore.NewShipmentRepository(db)
		shipmentUsecase := usecase.NewShipmentUseCase(shipmentRepository, shipmentTransaction)
		shipmentHandler := handler.NewShipmentHandler(shipmentUsecase)
		return &ShipmentContainer{
			Repository: shipmentRepository,
			Usecase:    shipmentUsecase,
			Handler:    shipmentHandler,
		}
	default:
		log.Print("default")
		shipmentTransaction := datastore.NewTransaction(db)
		shipmentRepository := datastore.NewShipmentRepository(db)
		shipmentUsecase := usecase.NewShipmentUseCase(shipmentRepository, shipmentTransaction)
		shipmentHandler := handler.NewShipmentHandler(shipmentUsecase)
		return &ShipmentContainer{
			Repository: shipmentRepository,
			Usecase:    shipmentUsecase,
			Handler:    shipmentHandler,
		}
	}
}
