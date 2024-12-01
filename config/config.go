package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	App struct {
		Name string `mapstructure:"name"`
		Port string `mapstructure:"port"`
	} `mapstructure:"app"`
	Database struct {
		Host string `mapstructure:"host"`
		// ENV = APP_DATABASE_HOST
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"database"`
	Jwt struct {
		SecretKey string `mapstructure:"secret_key"`
	} `mapstructure:"jwt"`
}

func LoadConfig() error {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	v := viper.New()

	v.SetEnvPrefix("APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	v.AutomaticEnv()

	if env == "development" {
		v.SetConfigName("dev")
		v.SetConfigType("yaml")

		v.AddConfigPath(".")
		v.AddConfigPath("./config")

		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return fmt.Errorf("error reading config file: %w", err)
			}
			log.Println("No dev config file found, using environment variables")
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return fmt.Errorf("unable to decode configuration: %w", err)
	}

	Cfg = &config

	fmt.Println(Cfg)

	return nil
}
