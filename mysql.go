package xman

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var InvalidMysqlConfig = errors.New("invalid mysql Config")

type Mysql struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	Dbname       string `yaml:"db"`
	Config       string `yaml:"config"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
	MaxOpenConns int    `yaml:"max-open-conns"`
	LogMode      bool   `yaml:"log"`
}

func (p *Mysql) isValid() bool {
	if p.Host == "" || p.Port == 0 || p.Username == "" || p.Password == "" {
		return false
	}

	return true
}

// 初始化数据库并产生数据库全局变量
func NewMysql(config Mysql) *gorm.DB {
	if !config.isValid() {
		panic(InvalidMysqlConfig)
	}
	dsn := fmt.Sprintf("%s:%s@(%s:%d)/%s?%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Dbname,
		config.Config,
	)
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(fmt.Sprintf("MySQL链接异常, %v", err))
	} else {
		db.DB().SetMaxIdleConns(config.MaxIdleConns)
		db.DB().SetMaxOpenConns(config.MaxOpenConns)
		db.LogMode(config.LogMode)
	}

	return db
}
