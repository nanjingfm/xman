package xman

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	oplogging "github.com/op/go-logging"
	"github.com/spf13/viper"
)

var (
	_db        *gorm.DB
	_cache     *redis.Client
	_config    Server
	_appConfig *viper.Viper
	_log       *oplogging.Logger
)

func Conf() *viper.Viper {
	return _appConfig
}

func LoadConf(i interface{}) error {
	c := Conf()
	if err := c.Unmarshal(i); err != nil {
		return err
	}

	return nil
}

func DB() *gorm.DB {
	return _db
}

func Cache() *redis.Client {
	return _cache
}

func sysConf() Server {
	return _config
}

func IsDev() bool {
	return sysConf().System.Env == "dev"
}

func LogError(args ...interface{}) {
	_log.Error()
}

func LogInfo(args ...interface{}) {
	_log.Info()
}

func LogDebug(args ...interface{}) {
	_log.Debug()
}
