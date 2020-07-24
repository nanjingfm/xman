package xman

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_initRedis(t *testing.T) {
	_config.Redis.Addr = ""
	assert.Panics(t, func() {
		initRedis()
	}, InvalidRedisConfig.Error())

	_config.Redis.Addr = "127.0.0.1:6379"
	initRedis()
	assert.NotNil(t, _cache)

	_cache = nil
	_config.Redis.Addr = "127.0.0.1:63"
	assert.Panics(t, func() {
		initRedis()
	}, "Redis连接异常")
}
