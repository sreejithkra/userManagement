package config

import (
	"log"

	"github.com/spf13/viper"
)

type Env struct {
	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
}

func EnvConfig() Env {
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	env := Env{
		DB_PORT:     viper.GetString("PORT"),
		DB_USER:     viper.GetString("DB_USER"),
		DB_PASSWORD: viper.GetString("DB_PASSWORD"),
		DB_HOST:     viper.GetString("DB_HOST"),
		DB_NAME:     viper.GetString("DB_NAME"),
	}

	return env
}
