package datastore

import (
	"context"
	"database/sql"

	"github.com/Tatsuemon/isucon9-qualify-shipment/domain/entity"

	"github.com/Tatsuemon/isucon9-qualify-shipment/domain/repository"

	"github.com/jmoiron/sqlx"
)

type shipmentRepository struct {
	conn *sqlx.DB
}

func NewShipmentRepository(conn *sqlx.DB) repository.ShipmentRepository {
	return &shipmentRepository{conn: conn}
}

func (r *shipmentRepository) FindByID(id string) (*entity.Shipment, error) {
	shipment := entity.Shipment{}
	if err := r.conn.Get(&shipment, "SELECT id, to_address, to_name, from_address, from_name, status, reserve_date_time, created_at FROM shipments WHERE id = ?", id); err != nil {
		return nil, err
	}

	return &shipment, nil
}

func (r *shipmentRepository) Store(ctx context.Context, s *entity.Shipment) (*entity.Shipment, error) {
	var tx interface {
		Prepare(query string) (*sql.Stmt, error)
	}
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.conn
	}
	stmt, err := tx.Prepare("INSERT INTO `shipments` (id, to_address, to_name, from_address, from_name, status, reserve_date_time, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = stmt.Exec(s.ID, s.ToAddress, s.ToName, s.FromAddress, s.FromName, s.Status, s.ReserveDateTime, s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *shipmentRepository) Update(ctx context.Context, s *entity.Shipment) (*entity.Shipment, error) {
	var tx interface {
		Prepare(query string) (*sql.Stmt, error)
	}
	tx, ok := GetTx(ctx)
	if !ok {
		tx = r.conn
	}
	stmt, err := tx.Prepare("UPDATE `shipments` SET to_address=?, to_name=?, from_address=?, from_name=?, status=?, reserve_date_time=?, created_at=? WHERE id=?")

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = stmt.Exec(s.ToAddress, s.ToName, s.FromAddress, s.FromName, s.Status, s.ReserveDateTime, s.CreatedAt, s.ID)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *shipmentRepository) Delete(id string) error {
	stmt, err := r.conn.Prepare("DELETE FROM `shipments` WHERE id=?")

	if err != nil {
		return err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}
	return nil
}
