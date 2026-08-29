package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	AppName           string
	AppEnv            string
	DBConnection      string
	TokenSymmetricKey string
	HttpUrl           string
	HttpPort          string
	DBHost            string
	DBPort            string
	DBdatabase        string
	DBUsername        string
	DBPassword        string
	JWTSecretKey      string
	JWTTokenDuration  time.Duration
}

func Load() Config {

	// ---------------------------------------------------------
	// 1. Configure Viper
	// ---------------------------------------------------------

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Look for config.yaml in these locations.
	viper.AddConfigPath(".")
	viper.AddConfigPath("/app")

	// Allow environment variables to be used as configuration.
	viper.AutomaticEnv()

	// ---------------------------------------------------------
	// 2. Explicitly map environment variables
	// ---------------------------------------------------------

	if err := viper.BindEnv("app_env", "APP_ENV"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("http_url", "HTTP_URL"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("http_port", "HTTP_PORT"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("db_host", "DB_HOST"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("db_port", "DB_PORT"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("db_database", "DB_DATABASE"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("db_username", "DB_USERNAME"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("db_password", "DB_PASSWORD"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("jwt_secret_key", "JWT_SECRET_KEY"); err != nil {
		log.Fatal(err)
	}

	if err := viper.BindEnv("jwt_token_duration", "JWT_TOKEN_DURATION"); err != nil {
		log.Fatal(err)
	}

	// ---------------------------------------------------------
	// 3. Read config.yaml
	// ---------------------------------------------------------

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Failed to read config file:", err)
		os.Exit(1)
	}

	// ---------------------------------------------------------
	// 4. Convert Viper configuration into Go Config struct
	// ---------------------------------------------------------

	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		fmt.Println("Failed to unmarshal configuration:", err)
		os.Exit(1)
	}

	// ---------------------------------------------------------
	// 5. Parse JWT token duration
	// ---------------------------------------------------------

	durationStr := viper.GetString("jwt_token_duration")

	if durationStr == "" {
		durationStr = "1h"
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Fatalf(
			"Invalid JWT_TOKEN_DURATION in config: %v",
			err,
		)
	}

	c.JWTTokenDuration = duration

	// ---------------------------------------------------------
	// 6. Set environment variables
	//
	// Existing parts of your application use os.Getenv().
	// Therefore we preserve compatibility with that code.
	// ---------------------------------------------------------

	os.Setenv("APP_ENV", c.AppEnv)
	os.Setenv("DB_CONNECTION", c.DBConnection)
	os.Setenv("TOKEN_SYMMETRIC_KEY", c.TokenSymmetricKey)

	os.Setenv("HTTP_URL", c.HttpUrl)
	os.Setenv("HTTP_PORT", c.HttpPort)

	os.Setenv("DB_HOST", c.DBHost)
	os.Setenv("DB_PORT", c.DBPort)
	os.Setenv("DB_DATABASE", c.DBdatabase)
	os.Setenv("DB_USERNAME", c.DBUsername)
	os.Setenv("DB_PASSWORD", c.DBPassword)

	os.Setenv("JWT_SECRET_KEY", c.JWTSecretKey)
	os.Setenv("JWT_TOKEN_DURATION", durationStr)

	return c
}
