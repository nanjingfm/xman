package xman

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_initMysql(t *testing.T)  {
	_config.Mysql = Mysql{}
	assert.Panics(t, func() {
		initMysql()
	}, InvalidMysqlConfig.Error())

	_config.Mysql = Mysql{
		Username:     "root",
		Password:     "123456",
		Host:         "127.0.0.1:3306",
		Dbname:       "mysql",
	}
	initMysql()
	assert.NotNil(t, _db)

	_db = nil
	_config.Mysql = Mysql{
		Username:     "root",
		Password:     "123456",
		Host:         "127.0.0.1:33",
		Dbname:       "mysql",
	}
	assert.Panics(t, func() {
		initMysql()
	}, "MySQL链接异常")
}
