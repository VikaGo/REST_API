package config

import (
	"encoding/json"
	"fmt"
	"github.com/kelseyhightower/envconfig"
	"log"
	"sync"
)

type Config struct {
	DbHost           string `envconfig:"DB_HOST"`
	DbPort           int    `envconfig:"DB_PORT"`
	DbUser           string `envconfig:"DB_USER"`
	DbPassword       string `envconfig:"DB_PASSWORD"`
	DbName           string `envconfig:"DB_NAME"`
	LogLevel         string `envconfig:"LOG_LEVEL"`
	PgURL            string `envconfig:"PG_URL"`
	PgMigrationsPath string `envconfig:"PG_MIGRATIONS_PATH"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment. Once.
func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Configuration:", string(configBytes))
	})
	return &config
}
