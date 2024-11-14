package config

import (
	logmgr "judgement/log"

	"github.com/spf13/viper"
)

type Config struct {
	Parts    []Users `yaml:"parts"`
	LogLevel string  `yaml:"logLevel,omitempty"`
}

type Users struct {
	Name                 string `yaml:"name"`
	ReplyRobotWebhookUrl string `yaml:"replyRobotWebhookUrl"`
}

func LoadConfig() {
	viper.SetConfigFile("/config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		logmgr.Log.Error("Error reading config file: %v", err)
	}
}
