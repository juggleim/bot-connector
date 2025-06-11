package configures

import (
	"os"

	"gopkg.in/yaml.v2"
)

type ConnectConfig struct {
	Port int `yaml:"port"`

	ApiKeySecret string `yaml:"apiKeySecret"`

	Log struct {
		LogPath string `yaml:"logPath"`
		LogName string `yaml:"logName"`
	} `ymal:"log"`

	Mysql struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Address  string `yaml:"address"`
		DbName   string `yaml:"name"`
	} `yaml:"mysql"`

	ImApiDomain string `yaml:"imApiDomain"`

	BotConnector struct {
		Domain string `yaml:"domain"`
	} `yaml:"botConnector"`
}

var Config ConnectConfig

func InitConfigures() error {
	cfBytes, err := os.ReadFile("conf/config.yml")
	if err != nil {
		return err
	}
	var conf ConnectConfig
	yaml.Unmarshal(cfBytes, &conf)
	Config = conf
	if Config.ApiKeySecret == "" {
		Config.ApiKeySecret = "1default-apikey1"
	}
	return nil
}
