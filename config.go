package xman

import (
	"fmt"
	"github.com/spf13/viper"
)

var _sysConfigFile = "config/sys.yaml"
var _appConfigFile = "config/app.yaml"

type Env string

const (
	EnvProd Env = "prod"
	EnvDev  Env = "dev"
)

type Server struct {
	DB     DB          `yaml:"mysql"`
	Redis  Redis       `yaml:"redis"`
	System System      `yaml:"system"`
	Log    LogConfig   `yaml:"log"`
	I18n   I18nOptions `yaml:"i18n"`
}

type System struct {
	Env        Env    `yaml:"env"`
	Addr       int    `yaml:"addr"`
	SigningKey string `yaml:"signing-key"`
}

func NewAppConfig(configPath string) *viper.Viper {
	if configPath == "" {
		configPath = _appConfigFile
	}
	vp := viper.New()
	vp.SetConfigFile(configPath)
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	return vp
}

// NewSysConfig 系统参数配置
func NewSysConfig(configPath string) Server {
	if configPath == "" {
		configPath = _appConfigFile
	}
	vp := viper.New()
	vp.SetConfigFile(configPath)
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	config := Server{}
	if err := vp.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("Fatal error unmarshal config file: %s \n", err))
	}
	return config
}
