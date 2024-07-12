package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	lop "github.com/samber/lo/parallel"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Get ENV
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASS", "postgres")
	viper.SetDefault("DB_NAME", "candahome")
	viper.SetDefault("SERVER_PORT", "8080")

	envBytes, err := assetsEnvBytes()
	if err != nil {
		log.Fatal(err.Error())
	}
	err = viper.ReadConfig(bytes.NewReader(envBytes))
	if err != nil {
		log.Fatal(err.Error())
	}

	// Grab files
	rowHTML, err := assetsTemplatesRowHtmlBytes()
	if err != nil {
		log.Fatal(err.Error())
	}
	viper.Set("ROW_HTML", rowHTML)

	pageHTML, err := assetsTemplatesPageHtmlBytes()
	if err != nil {
		log.Fatal(err.Error())
	}
	viper.Set("PAGE_HTML", pageHTML)

	// Connect to DB
	dsn := "host=" + viper.GetString("DB_HOST") +
		" user=" + viper.GetString("DB_USER") +
		" password=" + viper.GetString("DB_PASS") +
		" dbname=" + viper.GetString("DB_NAME") +
		" port=" + viper.GetString("DB_PORT")
	dbHandle, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	err = InitializeDB(dbHandle)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Build server tools
	serverContext := HTTPContext{
		dbHandle,
	}

	// Build router
	router := mux.NewRouter()
	router.HandleFunc("/sku", serverContext.CreateSKU).Methods(http.MethodPost)
	router.HandleFunc("/sku/{skuId}", serverContext.UpdateSKU).Methods(http.MethodPatch, http.MethodDelete)

	router.HandleFunc("/", serverContext.ReadSKUs).Methods(http.MethodGet)
	router.HandleFunc("/{location}", serverContext.ReadSKUs).Methods(http.MethodGet)

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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestBody.Location == "" {
		requestBody.Location = "pantry"
	} else {
		requestBody.Location = strings.ToLower(requestBody.Location)
	}

	err = CreateSKU(ctx.DB, requestBody)
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
	// Get tag from query
	vars := mux.Vars(r)
	location := vars["location"]
	if location == "" {
		location = "pantry"
	} else {
		location = strings.ToLower(location)
	}

	otherTags := lop.Map(r.URL.Query()["tags"], func(el string, _ int) string {
		return strings.ToLower(el)
	})

	skus, err := ReadSKUs(ctx.DB, append(otherTags, location))

	if err != nil {
		errObj := APIResponseBody{false, err.Error()}
		body, _ := json.Marshal(errObj)
		w.Header().Add("Content-Type", "application/json")
		http.Error(w, string(body), http.StatusInternalServerError)
	} else {
		html := BuildHTML(skus, location)

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
		err = DeleteSKU(ctx.DB, skuId)
	} else {
		// PATCH
		var requestBody SKUUpdateBody

		jsonDecoder := json.NewDecoder(r.Body)
		jsonDecoder.DisallowUnknownFields()
		err = jsonDecoder.Decode(&requestBody)
		if err == nil {
			err = UpdateSKU(ctx.DB, skuId, requestBody)
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
