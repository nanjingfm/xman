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
	DB     DB          `mapstructure:"db" yaml:"db"`
	Redis  Redis       `mapstructure:"redis" yaml:"redis"`
	System System      `mapstructure:"system" yaml:"system"`
	Log    LogConfig   `mapstructure:"log" yaml:"log"`
	I18n   I18nOptions `mapstructure:"i18n" yaml:"i18n"`
}

type System struct {
	Env        Env    `mapstructure:"env" yaml:"env"`
	Addr       int    `mapstructure:"addr" yaml:"addr"`
	SigningKey string `mapstructure:"signing-key" yaml:"signing-key"`
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
	vp.SetConfigType("yaml")
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
