// Utils for all HTML rendering functionality
package main

import (
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

// BuildHTML will render the SKU s in a simple HTML page
func BuildHTML(skus []SKU) (html string) {
	rows := ""
	for _, sku := range skus {
		rows += renderRow(sku)
	}

	return strings.Replace(viper.GetString("PAGE_HTML"), "%%ROWS%%", rows, 1)
}

func renderRow(sku SKU) (row string) {
	return strings.Replace(
		strings.Replace(
			strings.Replace(
				viper.GetString("ROW_HTML"),
				"%%SKU_QUANTITY%%",
				strconv.Itoa(int(sku.SkuQuantity)),
				-1,
			),
			"%%SKU_NAME%%",
			sku.SkuName,
			-1,
		),
		"%%SKU_ID%%",
		sku.ID,
		-1,
	)
}
