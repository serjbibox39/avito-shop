package storage

import (
	"avito-shop/pkg/storage/models"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func getMerchList(db *pgxpool.Pool) (map[string]models.Merch, error) {
	out := make(map[string]models.Merch)
	q := `SELECT id, name, price FROM merch`

	rows, err := db.Query(context.Background(), q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var m models.Merch
		err = rows.Scan(
			&m.ID,
			&m.Name,
			&m.Price,
		)
		if err != nil {
			return nil, err
		}
		out[m.Name] = m
	}
	return out, rows.Err()
}
