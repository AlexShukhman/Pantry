package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Grab files
	rowHTML, err := os.ReadFile("templates/row.html")
	if err != nil {
		log.Fatal(err.Error())
	}
	pageHTML, err := os.ReadFile("templates/page.html")
	if err != nil {
		log.Fatal(err.Error())
	}

	// Config
	viper.SetConfigFile(".env")

	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASS", "postgres")
	viper.SetDefault("DB_NAME", "candahome")

	viper.SetDefault("SERVER_PORT", "8080")

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	viper.Set("ROW_HTML", rowHTML)
	viper.Set("PAGE_HTML", pageHTML)

	dbURL := "postgres://" + viper.GetString("DB_USER") + ":" + viper.GetString("DB_PASS") +
		"@" + viper.GetString("DB_HOST") + ":" + viper.GetString("DB_PORT") +
		"/" + viper.GetString("DB_NAME")
	dbPool, err := pgxpool.New(context.Background(), dbURL)

	serverContext := HTTPContext{
		dbPool,
	}

	router := mux.NewRouter()
	router.HandleFunc("/sku", serverContext.CreateSKU).Methods(http.MethodPost)
	router.HandleFunc("/", serverContext.ReadSKUs).Methods(http.MethodGet)
	router.HandleFunc("/sku/{skuId}", serverContext.UpdateSKU).Methods(http.MethodPatch, http.MethodDelete)

	// Start server
	fmt.Printf("Starting server on port %d... \n", viper.GetInt("SERVER_PORT"))
	err = http.ListenAndServe(":"+viper.GetString("SERVER_PORT"), router)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (ctx *HTTPContext) CreateSKU(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var requestBody SKUCreateBody

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	err := jsonDecoder.Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = CreateSKU(ctx.DBPool, requestBody)
	if err != nil {
		errObj := APIResponseBody{false, err.Error()}
		body, _ := json.Marshal(errObj)
		http.Error(w, string(body), http.StatusInternalServerError)
	} else {
		bodyObj := APIResponseBody{true, ""}
		body, _ := json.Marshal(bodyObj)
		io.WriteString(w, string(body))
	}
}

func (ctx *HTTPContext) ReadSKUs(w http.ResponseWriter, r *http.Request) {
	skus, err := ReadSKUs(ctx.DBPool)

	if err != nil {
		errObj := APIResponseBody{false, err.Error()}
		body, _ := json.Marshal(errObj)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, string(body), http.StatusInternalServerError)
	} else {
		html := BuildHTML(skus)

		io.WriteString(w, html)
	}
}

func (ctx *HTTPContext) UpdateSKU(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	vars := mux.Vars(r)
	skuId := vars["skuId"]

	var err error
	if r.Method == http.MethodDelete {
		// DELETE
		err = DeleteSKU(ctx.DBPool, skuId)
	} else {
		// PATCH
		var requestBody SKUUpdateBody

		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()
		err = jsonDecoder.Decode(&requestBody)
		if err == nil {
			err = UpdateSKU(ctx.DBPool, skuId, requestBody)
		}
	}

	if err != nil {
		errObj := APIResponseBody{false, err.Error()}
		body, _ := json.Marshal(errObj)
		http.Error(w, string(body), http.StatusInternalServerError)
	} else {
		bodyObj := APIResponseBody{true, ""}
		body, _ := json.Marshal(bodyObj)
		io.WriteString(w, string(body))
	}
}
