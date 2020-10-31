package datastore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type DataBase interface {
	Close() error
}

type MysqlDB struct {
	DB *sqlx.DB
}

func NewMysqlDB(datasource string) (*MysqlDB, error) {
	db, err := sqlx.Open("mysql", datasource)
	if err != nil {
		return nil, fmt.Errorf("failed to open MySQL: %w", err)
	}

	db.SetMaxIdleConns(100) // idleコネクションの総数
	db.SetMaxOpenConns(100) // 全体のコネクションの総数

	// 接続確認
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping: %w", err)
	}

	return &MysqlDB{DB: db}, nil
}

func (m *MysqlDB) Close() error {
	return m.DB.Close()
}
