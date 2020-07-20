package xman

import (
	"errors"
	"github.com/go-redis/redis"
)

var InvalidRedisConfig = errors.New("invalid redis Config")

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`
}

func (p *Redis) isValid() bool {
	return p.Addr != ""
}

func initRedis() {
	redisCfg := sysConf().Redis
	if !redisCfg.isValid() {
		panic(InvalidRedisConfig)
	}
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
		LogError(err)
	} else {
		LogInfo("redis connect ping response:", pong)
		_cache = client
	}
}
