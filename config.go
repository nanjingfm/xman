package xman

import (
	"fmt"

	"github.com/spf13/viper"
)

var _sysConfigFile = "config/sys.yaml"
var _appConfigFile = "config/app.yaml"

type Server struct {
	Mysql   Mysql       `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis   Redis       `mapstructure:"redis" json:"redis" yaml:"redis"`
	System  System      `mapstructure:"system" json:"system" yaml:"system"`
	Captcha Captcha     `mapstructure:"captcha" json:"captcha" yaml:"captcha"`
	Log     Log         `mapstructure:"log" json:"log" yaml:"log"`
	I18n    I18nOptions `mapstructure:"i18n" json:"i18n" yaml:"i18n"`
}

type System struct {
	UseMultipoint bool   `mapstructure:"use-multipoint" json:"useMultipoint" yaml:"use-multipoint"`
	Env           string `mapstructure:"env" json:"env" yaml:"env"`
	Addr          int    `mapstructure:"addr" json:"addr" yaml:"addr"`
	DbType        string `mapstructure:"db-type" json:"dbType" yaml:"db-type"`
	SigningKey    string `mapstructure:"signing-key" json:"signingKey" yaml:"signing-key"`
}

func initConfig() {
	initSysConfig()
	initAppConfig()
}

func initAppConfig() {
	vp := viper.New()
	vp.SetConfigFile(_appConfigFile)
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	//vp.Debug()
	//vp.WatchConfig()
	_appConfig = vp
}

// initSysConfig 系统参数配置
func initSysConfig() {
	vp := viper.New()
	vp.SetConfigFile(_sysConfigFile)
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err := vp.Unmarshal(&_config); err != nil {
		fmt.Println(err)
	}
}
