package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Configuration struct {
	ENV []Settings `mapstructure:"env"`
}

type Settings struct {
	NAME  string
	VALUE string
}

func InitConfig(environment string) {

	viper.SetConfigType("yaml")

	switch environment {
	case "local":
		viper.SetConfigName("values.local")
	default:
		log.Fatalf("configuration for %s not found", environment)
	}

	viper.AddConfigPath("./config")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	var env Configuration
	viper.Unmarshal(&env)
	for _, value := range env.ENV {
		if _, ok := os.LookupEnv(value.NAME); !ok {
			os.Setenv(value.NAME, value.VALUE)
		}
	}
}
