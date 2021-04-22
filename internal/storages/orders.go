package storages

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/peakle/recycle/internal"
	idgen "github.com/wakeapp/go-id-generator"
)

func GetList(ctx context.Context, m *internal.SQLManager) ([]*Order, error) {
	var rows *sql.Rows
	var conn = m.GetConnection()
	var stmt, err = conn.PrepareContext(ctx, "SELECT id, address, maxSize, eventAt FROM Orders")
	if err != nil {
		return nil, fmt.Errorf("on GetList: on Prepare: %s", err)
	}
	defer stmt.Close()

	rows, err = stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("on GetList: on Query: %s", err)
	}
	defer rows.Close()

	var orders = make([]*Order, 0, 10)
	var order Order
	for rows.Next() {
		err = rows.Scan(
			&order.Id,
			&order.Address,
			&order.MaxSize,
			&order.EventAt,
		)

		if err != nil {
			return nil, fmt.Errorf("on GetList: on Scan: %s", err)
		}

		orders = append(orders, &order)
	}

	return orders, nil
}

func Info(ctx context.Context, m *internal.SQLManager, orderId string) (*Order, error) {
	var rows *sql.Rows
	var conn = m.GetConnection()
	var stmt, err = conn.PrepareContext(ctx,
		`
		SELECT
			id,
			address,
			maxSize,
			eventAt,
			count(uo.user_id)
		FROM Orders o
		JOIN OrdersUsers ou on o.id = ou.order_id
		WHERE id = ?
		LIMIT 1
	`)
	if err != nil {
		return nil, fmt.Errorf("on GetList: on Prepare: %s", err)
	}
	defer stmt.Close()

	rows, err = stmt.QueryContext(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("on Info: on Query: %s", err)
	}
	defer rows.Close()

	var order Order
	for rows.Next() {
		err = rows.Scan(
			&order.Id,
			&order.Address,
			&order.MaxSize,
			&order.EventAt,
			&order.CurrentSize,
		)

		if err != nil {
			return nil, fmt.Errorf("on Info: on Scan: %s", err)
		}
	}

	return &order, nil
}

func Subscribe(ctx context.Context, m *internal.SQLManager, orderId, userId string) (bool, error) {
	var conn = m.GetConnection()
	var stmt, err = conn.PrepareContext(ctx, "INSERT INTO OrdersUsers (order_id, user_id) VALUES(?,?)")
	if err != nil {
		return false, fmt.Errorf("on Subscribe: on Prepare: %s", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, orderId, userId)
	if err != nil {
		return false, fmt.Errorf("on Subscribe: on Exec: %s", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("on Subscribe: on RowsAffected: %s", err)
	}

	if count == 0 {
		return false, fmt.Errorf("on Subscribe: zero rows affected")
	}

	return true, nil
}

func Create(ctx context.Context, m *internal.SQLManager, userId, address, maxSize, eventAt string) (string, error) {
	var err error
	var stmt *sql.Stmt
	var tx *sql.Tx
	var res sql.Result
	var count int64

	var conn = m.GetConnection()
	tx, err = conn.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return "", fmt.Errorf("on Create: on start transaction: %s", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}

		tx.Commit()
	}()

	stmt, err = tx.PrepareContext(ctx, "INSERT INTO Orders (id, address, maxSize, eventAt, createdAt, updatedAt) VALUES(?,?,?,?,?,?)")
	if err != nil {
		return "", fmt.Errorf("on Create: on Prepare: %s", err)
	}
	defer func(s *sql.Stmt) {
		_ = s.Close()
	}(stmt)

	id := strings.TrimLeft(idgen.Id(), "0")

	var t = time.Now().Format("2006-01-02 15:04:05")
	res, err = stmt.ExecContext(ctx, id, address, maxSize, eventAt, t, t)
	if err != nil {
		return "", fmt.Errorf("on Create: on Exec: %s", err)
	}

	count, err = res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("on Create: on RowsAffected: %s", err)
	}

	if count == 0 {
		return "", fmt.Errorf("zero rows affected")
	}

	stmt, err = tx.PrepareContext(ctx, "INSERT INTO OrdersUsers (user_id, order_id) VALUES(?,?)")
	if err != nil {
		return "", fmt.Errorf("on Create: on Prepare: %s", err)
	}
	defer func(s *sql.Stmt) {
		_ = s.Close()
	}(stmt)

	res, err = stmt.ExecContext(ctx, userId, id)
	if err != nil {
		return "", fmt.Errorf("on Create: on Exec: %s", err)
	}

	count, err = res.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("on Create: on RowsAffected: %s", err)
	}

	if count == 0 {
		return "", fmt.Errorf("zero rows affected")
	}

	return id, nil
}
