package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config Global Configuration
var CONFIG *Configurations

func init() {
	var err error
	var configuration Configurations

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}

	err = viper.Unmarshal(&configuration)
	if err != nil {
		fmt.Printf("Unable to decode into struct, %v", err)
	}

	CONFIG = &configuration
}