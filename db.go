package xman

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var InvalidMysqlConfig = errors.New("invalid mysql Config")

func InsertIgnore() (string, interface{}) {
	return "gorm:insert_modifier", "IGNORE"
}

func InsertOnDuplicate(updateRule string) (string, interface{}) {
	return "gorm:insert_modifier", "ON DUPLICATE KEY UPDATE " + updateRule
}

type DB struct {
	DbType       string `mapstructure:"db_type" yaml:"db_type"`
	Dsn          string `mapstructure:"dsn" yaml:"dsn"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log" yaml:"log"`
}

func (p *DB) isValid() bool {
	if p.Dsn == "" {
		return false
	}

	if p.DbType != "mysql" && p.DbType != "mssql" {
		return false
	}

	return true
}

// 初始化数据库并产生数据库全局变量
func NewDB(config DB) *gorm.DB {
	if !config.isValid() {
		panic(InvalidMysqlConfig)
	}
	db, err := gorm.Open(config.DbType, config.Dsn)
	if err != nil {
		panic(fmt.Sprintf("MySQL链接异常, %v", err))
	} else {
		db.DB().SetMaxIdleConns(config.MaxIdleConns)
		db.DB().SetMaxOpenConns(config.MaxOpenConns)
		db.LogMode(config.LogMode)
	}

	return db
}
