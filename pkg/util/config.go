package util

import (
	"strings"

	"github.com/spf13/viper"
)

type EnvVar struct {
	Port      string `mapstructure:"SMC_AGENT_PORT"`
	ShowImage bool   `mapstructure:"SMC_SHOW_IMAGE"`
	ShowLaser bool   `mapstructure:"SMC_SHOW_LASER"`
	ShowNhso  bool   `mapstructure:"SMC_SHOW_NHSO"`
	ApiKey    string `mapstructure:"API_KEY"`
}
type Config struct {
	EnvVar
}

func LoadConfig(path string) (env EnvVar, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&env)
	return
}

func LoadConfigFromString(configStr string) (EnvVar, error) {
	viper.SetConfigType("env")
	err := viper.ReadConfig(strings.NewReader(configStr))
	if err != nil {
		return EnvVar{}, err
	}

	var config EnvVar
	err = viper.Unmarshal(&config)
	if err != nil {
		return EnvVar{}, err
	}

	return config, nil
}
