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

func Info(ctx context.Context, m *internal.SQLManager) (*Order, error) {
	var rows *sql.Rows
	var conn = m.GetConnection()
	var stmt, err = conn.PrepareContext(ctx, "SELECT id, address, maxSize, eventAt FROM Orders LIMIT 1")
	if err != nil {
		return nil, fmt.Errorf("on GetList: on Prepare: %s", err)
	}
	defer stmt.Close()

	rows, err = stmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("on GetList: on Query: %s", err)
	}
	defer rows.Close()

	var order Order
	for rows.Next() {
		err = rows.Scan(
			&order.Id,
			&order.Address,
			&order.MaxSize,
			&order.EventAt,
		)

		if err != nil {
			return nil, fmt.Errorf("on Info: on Scan: %s", err)
		}
	}

	return &order, nil
}
