package main

import (
	"fmt"
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

	// Start server
	fmt.Printf("Starting server on port " + viper.GetString("SERVER_PORT") + "\n")
}
