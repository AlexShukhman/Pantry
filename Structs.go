package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HTTPContext struct {
	DB *gorm.DB
}

type APIResponseBody struct {
	Success     bool   `json:"success"`
	ErrorString string `json:"errorString"`
}

type SKU struct {
	ID          uuid.UUID `json:"id" gorm:"primaryKey;not null;type:uuid;default:gen_random_uuid()"`
	SkuName     string    `json:"skuName" gorm:"index;column:sku_name;not null;unique;type:text"`
	SkuQuantity int16     `json:"skuQuantity" gorm:"column:sku_quantity;not null;type:smallint"`
	SKUTags     []SKUTag  `gorm:"foreignKey:sku_id;OnDelete:CASCADE"`
}

type SKUTag struct {
	SKUId uuid.UUID `json:"skuId" gorm:"primaryKey;autoIncrement:false;column:sku_id;type:uuid;not null"`
	TagId string    `json:"tagId" gorm:"primaryKey;autoIncrement:false;column:tag_id;type:varchar(255);not null"`
}

type Tag struct {
	ID      string   `json:"id" gorm:"primaryKey;not null;type:varchar(255)"`
	SKUTags []SKUTag `gorm:"foreignKey:tag_id;OnDelete:CASCADE"`
}

type SKUCreateBody struct {
	SkuName     string   `json:"name"`
	SkuQuantity int16    `json:"quantity"`
	Location    string   `json:"location"`
	Tags        []string `json:"tags"`
}

type SKUUpdateBody struct {
	AdditionalQuantity int16 `json:"additionalQuantity"`
}
