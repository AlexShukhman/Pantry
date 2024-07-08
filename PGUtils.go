// Utils for all Postgres-related functionality
package main

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CreateSKU will add a SKU to the PantrySKUs table
func CreateSKU(dbPool *pgxpool.Pool, sku SKUCreateBody) (err error) {
	query := `INSERT INTO pantryskus (id, sku_name, sku_quantity) VALUES ($1, $2, $3)`

	_, err = dbPool.Exec(
		context.Background(),
		query,
		uuid.New().String(),
		sku.SkuName,
		sku.SkuQuantity,
	)

	return err
}

// ReadSKUs will get all SKU s in PantrySKUs table
func ReadSKUs(dbPool *pgxpool.Pool) (renderedRows []SKU, err error) {
	query := `SELECT * FROM pantryskus ORDER BY sku_name;`
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

// UpdateSKU will update the SKU according to strict allowed updates
func UpdateSKU(dbPool *pgxpool.Pool, skuId string, update SKUUpdateBody) (err error) {
	query := `UPDATE pantryskus SET sku_quantity = sku_quantity + $1 WHERE id = $2`

	_, err = dbPool.Exec(context.Background(), query, update.AdditionalQuantity, skuId)

	return err
}

// DeleteSKU will delete the SKU by skuId
func DeleteSKU(dbPool *pgxpool.Pool, skuId string) (err error) {
	query := `DELETE FROM pantryskus WHERE id = $1`

	_, err = dbPool.Exec(context.Background(), query, skuId)

	return err
}
