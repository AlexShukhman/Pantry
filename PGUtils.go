// Utils for all Postgres-related functionality
package main

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// GetAllSKUs will get all SKU s in PantrySKUs table
func GetAllSKUs(dbPool *pgxpool.Pool) (renderedRows []SKU, err error) {
	query := `SELECT * FROM pantryskus;`
	rows, err := dbPool.Query(context.Background(), query)
	if err != nil {
		return
	}

	skus := []SKU{}
	for rows.Next() {
		sku := SKU{}
		err := rows.Scan(
			&sku.ID,
			&sku.SkuName,
			&sku.SkuQuantity,
		)
		if err != nil {
			return renderedRows, err
		}

		skus = append(skus, sku)
	}

	return skus, nil
}

// InsertSKU will add a SKU to the PantrySKUs table
func InsertSKU(dbPool *pgxpool.Pool, sku SKU) (err error) {
	query := `INSERT INTO pantryskus (id, sku_name, sku_quantity) VALUES (@id, @sku_name, @sku_quantity)`

	args := pgx.NamedArgs{
		"id":           sku.ID,
		"sku_name":     sku.SkuName,
		"sku_quantity": sku.SkuQuantity,
	}

	_, err = dbPool.Exec(context.Background(), query, args)

	return err
}
