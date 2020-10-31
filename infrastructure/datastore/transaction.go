package datastore

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

var txKey = struct{}{}

type Transaction interface {
	DoInTx(context.Context, func(context.Context) (interface{}, error)) (interface{}, error)
}

type tx struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) Transaction {
	return &tx{db: db}
}

func (t *tx) DoInTx(ctx context.Context, f func(ctx context.Context) (interface{}, error)) (interface{}, error) {
	tx, err := t.db.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	// ここでctxへトランザクションオブジェクトを放り込む
	ctx = context.WithValue(ctx, &txKey, tx)

	// トランザクションの対象処理へコンテキストを引継ぎ
	v, err := f(ctx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		// エラーならロールバック
		tx.Rollback()
		return nil, err
	}

	return v, nil
}

// context.Contextからトランザクションを取得する関数
func GetTx(ctx context.Context) (*sqlx.Tx, bool) {
	tx, ok := ctx.Value(&txKey).(*sqlx.Tx)
	return tx, ok
}
