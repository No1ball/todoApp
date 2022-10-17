package config

import "github.com/spf13/viper"

func InitConfig() error {
	viper.AddConfigPath("../todoApp/")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
