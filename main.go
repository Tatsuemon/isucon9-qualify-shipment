package main

import (
	"log"
	"os"

	"github.com/Tatsuemon/isucon9-qualify-shipment/config"
	di "github.com/Tatsuemon/isucon9-qualify-shipment/di/containers"
	"github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/datastore"
	server "github.com/Tatsuemon/isucon9-qualify-shipment/infrastructure/web"
)

func main() {
	const port = 8080

	db, err := datastore.NewMysqlDB(config.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal()
		}
	}()

	shipmentController := di.NewShipmentContainer(os.Getenv("ENV"), db.DB)
	s := server.NewServer()
	s.Init(shipmentController.Handler)
	s.Run(port)
}
