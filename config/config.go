package config

import (
	"github.com/spf13/viper"
	logmgr "judgement/config/log"
)

var AllConfig *Config
var (
	// CaptchaWebhookUrl 验证码提醒机器人地址

	CaptchaWebhookUrl string

	//客服对象列表

	Parts []Part

	// RemindWebhookUrl 发送提醒url

	RemindWebhookUrl string
)

type Config struct {
	LogLevel            string `yaml:"logLevel,omitempty"`
	CaptchaWebhookUrl   string `yaml:"captchaWebhookUrl"`
	RemindWebhookUrl    string `yaml:"remindWebhookUrl"`
	newTicketWebhookUrl string `yaml:"newTicketWebhookUrl"`
	Parts               []Part `yaml:"parts"`
}

type Part struct {
	Name            string `yaml:"name"`
	Phone           string `yaml:"phone"`
	ReplyWebhookUrl string `yaml:"replyWebhookUrl"`
}

func NewConfig(configPath string) *Config {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		logmgr.Log.Error("Error reading config file: %v", err)
		panic(err)
	}
	err = viper.Unmarshal(&AllConfig)
	if err != nil {
		logmgr.Log.Error("Error Unmarshal config file: %v", err)
		return nil
	}
	return AllConfig

}

func (a *Config) LoadConfig() {
	logmgr.LogLevel = a.LogLevel
	CaptchaWebhookUrl = a.CaptchaWebhookUrl
	Parts = a.Parts
	RemindWebhookUrl = a.RemindWebhookUrl
}
