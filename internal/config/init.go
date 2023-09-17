package config

import (
	"github.com/spf13/viper"
)

// init initialize default config params
func init() {
	// environment - could be "local", "prod", "dev"
	viper.SetDefault("env", "dev")

	viper.SetDefault("db.addr", "postgres://postgres:postgres@localhost:5432/pets?sslmode=disable")
	viper.SetDefault("db.driver", "postgres")

	viper.SetDefault("http.tcp", "0.0.0.0:8000")
}
