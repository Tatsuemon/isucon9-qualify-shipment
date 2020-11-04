package datastore

import (
	"context"
	"database/sql"
	"time"

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

type shipmentDTO struct {
	ID              string         `db:"id"`
	ToAddress       string         `db:"to_address"`
	ToName          string         `db:"to_name"`
	FromAddress     string         `db:"from_address"`
	FromName        string         `db:"from_name"`
	Status          string         `db:"status"`
	ReserveDateTime time.Time      `db:"reserve_date_time"`
	DoneDateTime    sql.NullTime   `db:"done_date_time"`
	QRMD5           sql.NullString `db:"qrmd5"`
	CreatedAt       time.Time      `db:"created_at"`
}

func (r *shipmentRepository) FindByID(id string) (*entity.Shipment, error) {
	dto := shipmentDTO{}
	if err := r.conn.Get(&dto, "SELECT id, to_address, to_name, from_address, from_name, status, reserve_date_time, done_date_time, qrmd5, created_at FROM shipments WHERE id = ?", id); err != nil {
		return nil, err
	}
	shipment := entity.Shipment{
		ID:              dto.ID,
		ToAddress:       dto.ToAddress,
		ToName:          dto.ToName,
		FromAddress:     dto.FromAddress,
		FromName:        dto.FromName,
		Status:          dto.Status,
		ReserveDateTime: dto.ReserveDateTime,
		CreatedAt:       dto.CreatedAt,
	}
	if dto.DoneDateTime.Valid {
		shipment.DoneDateTime = dto.DoneDateTime.Time
	}
	if dto.QRMD5.Valid {
		shipment.QRMD5 = dto.QRMD5.String
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
	stmt, err := tx.Prepare("INSERT INTO `shipments` (id, to_address, to_name, from_address, from_name, status, reserve_date_time, done_date_time, qrmd5, created_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	// time.Timeのゼロ値がnilでないため
	var doneDateTime interface{}
	if s.DoneDateTime.IsZero() {
		doneDateTime = nil
	} else {
		doneDateTime = s.DoneDateTime
	}

	_, err = stmt.Exec(s.ID, s.ToAddress, s.ToName, s.FromAddress, s.FromName, s.Status, s.ReserveDateTime, doneDateTime, s.QRMD5, s.CreatedAt)
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
	stmt, err := tx.Prepare("UPDATE `shipments` SET to_address=?, to_name=?, from_address=?, from_name=?, status=?, reserve_date_time=?, done_date_time=?, qrmd5=?, created_at=? WHERE id=?")

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := stmt.Close(); closeErr != nil {
			err = closeErr
		}
	}()

	// time.Timeのゼロ値がnilでないため
	var doneDateTime interface{}
	if s.DoneDateTime.IsZero() {
		doneDateTime = nil
	} else {
		doneDateTime = s.DoneDateTime
	}

	_, err = stmt.Exec(s.ToAddress, s.ToName, s.FromAddress, s.FromName, s.Status, s.ReserveDateTime, doneDateTime, s.QRMD5, s.CreatedAt, s.ID)
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
