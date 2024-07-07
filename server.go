package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx"
	"github.com/spf13/viper"
	"log"
)

func main() {
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

	dburl := "postgres://" + viper.GetString("DB_USER") + ":" + viper.GetString("DB_PASS") +
		"@" + viper.GetString("DB_HOST") + ":" + viper.GetString("DB_PORT") +
		"/" + viper.GetString("DB_NAME")
	conn, err := pgx.Connect(context.Background(), dburl)

	// Start server
	fmt.Printf("Starting server on port " + viper.GetString("SERVER_PORT") + "\n")
}
