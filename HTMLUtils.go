// Utils for all HTML rendering functionality
package main

import (
	"strconv"
	"strings"
)

var rowHTML string
var pageHTML string

// BuildHTML will render the SKU s in a simple HTML page
func BuildHTML(skus []SKU) (html string) {
	rows := ""
	for _, sku := range skus {
		rows += renderRow(sku)
	}

	return strings.Replace(pageHTML, "%%ROWS%%", rows, 1)
}

func renderRow(sku SKU) (row string) {
	return strings.Replace(
		strings.Replace(
			strings.Replace(
				rowHTML,
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
