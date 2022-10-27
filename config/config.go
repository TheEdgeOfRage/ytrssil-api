package config

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

type DB struct {
	DBURI string `long:"db-uri" env:"DB_URI"`
}

// Gin contains configuration for the gin framework
type Gin struct {
	Port int `long:"port" env:"PORT" default:"8080"`
}

// Config ties together all configs
type Config struct {
	DB  DB
	Gin Gin
}

func getenvOrDefault(key string, defaultValue string) string {
	value, found := os.LookupEnv(key)
	if found {
		return value
	}

	return defaultValue
}

// Parse parses all the supplied configurations and returns
func Parse() (Config, error) {
	var config Config
	parser := flags.NewParser(&config, flags.Default)
	_, err := parser.Parse()
	return config, err
}

// TestConfig returns a mostly hardcoded configuration used for running tests
func TestConfig() Config {
	dbURI := getenvOrDefault("DB_URI", "postgres://ytrssil:ytrssil@postgres:5432/ytrssil")

	gin := Gin{
		Port: 8080,
	}
	db := DB{
		DBURI: dbURI,
	}
	config := Config{
		Gin: gin,
		DB:  db,
	}

	return config
}
