package main

import "github.com/jackc/pgx/v5/pgxpool"

type HTTPContext struct {
	DBPool *pgxpool.Pool
}

type APIResponseBody struct {
	Success     bool   `json:"success"`
	ErrorString string `json:"errorString"`
}

type SKU struct {
	ID          string `json:"id"`
	SkuName     string `json:"skuName"`
	SkuQuantity int16  `json:"skuQuantity"`
}

type SKUCreateBody struct {
	SkuName     string `json:"name"`
	SkuQuantity int16  `json:"quantity"`
}

type SKUUpdateBody struct {
	AdditionalQuantity int16 `json:"additionalQuantity"`
}
