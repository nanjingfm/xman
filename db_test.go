package xman

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_initMysql(t *testing.T) {
	config := DB{}
	assert.Panics(t, func() {
		NewDB(config)
	}, InvalidMysqlConfig.Error())

	config1 := DB{
		DbType: "mysql",
		Dsn:    "root:123456@(127.0.0.1:3306)/mysql?charset=utf8&parseTime=True&loc=Local",
	}
	db := NewDB(config1)
	defer db.Close()
	assert.NotNil(t, db)

	config2 := DB{
		Dsn:    "root:123456@(127.0.0.1:33060)/mysql?charset=utf8&parseTime=True&loc=Local",
		DbType: "mysql",
	}
	assert.Panics(t, func() {
		db := NewDB(config2)
		defer db.Close()
	}, "MySQL链接异常")
}
