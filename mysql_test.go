package xman

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_initMysql(t *testing.T) {
	config := Mysql{}
	assert.Panics(t, func() {
		NewMysql(config)
	}, InvalidMysqlConfig.Error())

	config1 := Mysql{
		Username: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     3306,
		Dbname:   "mysql",
	}
	db := NewMysql(config1)
	defer db.Close()
	assert.NotNil(t, db)

	config2 := Mysql{
		Username: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     33060,
		Dbname:   "mysql",
	}
	assert.Panics(t, func() {
		db := NewMysql(config2)
		defer db.Close()
	}, "MySQL链接异常")
}
