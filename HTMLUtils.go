// Utils for all HTML rendering functionality
package main

import (
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strconv"
	"strings"
)

// BuildHTML will render the SKU s in a simple HTML page
func BuildHTML(skus []SKU, locTag string) (html string) {
	rows := ""
	skuOptions := ""
	for _, sku := range skus {
		rows += renderRow(sku)
		skuOptions += "<option>" + sku.SkuName + "</option>"
	}

	return strings.Replace(
		strings.Replace(
			strings.Replace(
				viper.GetString("PAGE_HTML"),
				"%%LOC_TAG%%",
				cases.Title(language.AmericanEnglish).String(locTag),
				1,
			),
			"%%SKU_OPTIONS%%",
			skuOptions,
			1,
		),
		"%%ROWS%%",
		rows,
		1,
	)
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
		sku.ID.String(),
		-1,
	)
}
