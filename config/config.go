package config

import (
	logmgr "judgement/log"

	"github.com/spf13/viper"
)

type Config struct {
	Parts            []Users `yaml:"parts"`
	RemindWebhookUrl string  `yaml:"remindWebhookUrl"`
	LogLevel         string  `yaml:"logLevel,omitempty"`
}

type Users struct {
	Name            string `yaml:"name"`
	Phone           string `yaml:"phone"`
	ReplyWebhookUrl string `yaml:"replyWebhookUrl"`
}

func LoadConfig() {
	viper.SetConfigFile("/config.yaml")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		logmgr.Log.Error("Error reading config file: %v", err)
		panic(err)
	}
}
