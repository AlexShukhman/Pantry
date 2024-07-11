// Utils for all Postgres-related functionality
package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// InitializeDB will add all required tables to the DB if it's missing them
func InitializeDB(db *gorm.DB) (err error) {
	return db.AutoMigrate(&SKU{}, &Tag{}, &SKUTag{})
}

// CreateSKU will add a SKU to the skus table
func CreateSKU(db *gorm.DB, sku SKUCreateBody) (err error) {
	newSKU := SKU{
		SkuName:     sku.SkuName,
		SkuQuantity: sku.SkuQuantity,
	}

	err = db.Select("sku_name", "sku_quantity").Create(newSKU).Error

	return err
}

// ReadSKUs will get all SKU s in skus table tagged "pantry"
func ReadSKUs(db *gorm.DB, tags []string) (renderedRows []SKU, err error) {
	// query := `SELECT DISTINCT skus.sku_name as sku_name, skus.id as id, skus.sku_quantity as sku_quantity FROM sku_tags RIGHT JOIN skus ON sku_tags.sku = skus.id WHERE sku_tags.tag IN $1 ORDER BY skus.sku_name;`

	var skus []SKU
	result := db.
		Joins("left join sku_tags on skus.id = sku_tags.sku_id").
		Where("sku_tags.tag_id IN ?", tags).
		Distinct(
			"skus.sku_name as sku_name",
			"skus.id as id",
			"skus.sku_quantity as sku_quantity",
		).
		Order("skus.sku_name").
		Find(&skus)

	if result.Error != nil {
		return renderedRows, err
	}
	return skus, nil
}

// UpdateSKU will update the SKU according to strict allowed updates
func UpdateSKU(db *gorm.DB, skuId string, update SKUUpdateBody) (err error) {
	// query := `UPDATE skus SET sku_quantity = sku_quantity + $1 WHERE id = $2`
	skuUUID, err := uuid.Parse(skuId)
	if err != nil {
		return
	}

	sku := SKU{
		ID: skuUUID,
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if update.AdditionalQuantity != 0 {
			txErr := tx.Model(sku).UpdateColumn("sku_quantity", gorm.Expr("sku_quantity + ?", update.AdditionalQuantity)).Error

			if txErr != nil {
				return txErr
			}
		}

		return nil
	})
}

// DeleteSKU will delete the SKU by skuId
func DeleteSKU(db *gorm.DB, skuId string) (err error) {
	// query := `DELETE FROM skus WHERE id = $1`
	skuUUID, err := uuid.Parse(skuId)
	if err != nil {
		return
	}

	return db.Delete(&SKU{}, skuUUID).Error
}
