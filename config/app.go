package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	AppEnv       string `mapstructure:"app_env"`
	DbConnection string `mapstructure:"db_connection"`
	DbHost       string `mapstructure:"db_host"`
	DbPort       string `mapstructure:"db_port"`
	DbDatabase   string `mapstructure:"db_database"`
	DbUsername   string `mapstructure:"db_username"`
	DbPassword   string `mapstructure:"db_password"`
	AppEnvIsProd bool
}

const APP_ENV_PROD = "prod"
const APP_ENV_DEV = "dev"

var Conf Config

func LoadConfig() {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln(err)
	}
	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatalln(err)
	}
	Conf.AppEnvIsProd = (Conf.AppEnv == APP_ENV_PROD)

}
