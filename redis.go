package xman

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
)

var InvalidRedisConfig = errors.New("invalid redis Config")

type Redis struct {
	Host     string `mapstructure:"addr" yaml:"addr"`
	Port     string `mapstructure:"port" yaml:"port"`
	Password string `mapstructure:"password" yaml:"password"`
	DB       int    `mapstructure:"db" yaml:"db"`
}

func (p *Redis) isValid() bool {
	return p.Host != ""
}

func NewRedis(config Redis) *redis.Client {
	if !config.isValid() {
		panic(InvalidRedisConfig)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     config.Host + ":" + config.Port,
		Password: config.Password, // no password set
		DB:       config.DB,       // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("Redis连接异常, %v", err))
	}
	return client
}
