package main

import "github.com/jackc/pgx/v5/pgxpool"

type HTTPContext struct {
	DBPool *pgxpool.Pool
}

type SKU struct {
	ID          string `json:"id"`
	SkuName     string `json:"skuName"`
	SkuQuantity int16  `json:"skuQuantity"`
}
