package xman

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_initRedis(t *testing.T) {
	config := Redis{}
	assert.Panics(t, func() {
		NewRedis(config)
	}, InvalidRedisConfig.Error())

	config2 := Redis{Host: "127.0.0.1", Port: "6379"}
	NewRedis(config2)
	assert.NotNil(t, config2)

	config3 := Redis{Host: "127.0.0.1", Port: "63719"}
	assert.Panics(t, func() {
		NewRedis(config3)
	}, "Redis连接异常")
}
