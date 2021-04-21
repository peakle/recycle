package storages

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/peakle/recycle/internal"
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

func Subscribe(ctx context.Context, m *internal.SQLManager, orderId string) (bool, error) {
	var conn = m.GetConnection()
	var stmt, err = conn.PrepareContext(ctx, "")
	if err != nil {
		return false, fmt.Errorf("on Subscribe: on Prepare: %s", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx)
	if err != nil {
		return false, fmt.Errorf("on Subscribe: on Exec: %s", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("on Subscribe: on RowsAffected: %s", err)
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}
