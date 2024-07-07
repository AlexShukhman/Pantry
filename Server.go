package main

import (
	"context"
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

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err.Error())
	}

	dbURL := "postgres://" + viper.GetString("DB_USER") + ":" + viper.GetString("DB_PASS") +
		"@" + viper.GetString("DB_HOST") + ":" + viper.GetString("DB_PORT") +
		"/" + viper.GetString("DB_NAME")
	dbPool, err := pgxpool.New(context.Background(), dbURL)

	serverContext := HTTPContext{
		dbPool,
	}

	router := mux.NewRouter()
	router.HandleFunc("/", serverContext.GetSKUs).Methods("GET")

	// Start server
	fmt.Printf("Starting server on port %d... \n", viper.GetInt("SERVER_PORT"))
	err = http.ListenAndServe(":"+viper.GetString("SERVER_PORT"), router)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (ctx *HTTPContext) GetSKUs(w http.ResponseWriter, r *http.Request) {
	skus, err := GetAllSKUs(ctx.DBPool)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	html := BuildHTML(skus)

	io.WriteString(w, html)
}
