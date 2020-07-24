package xman

import (
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var InvalidMysqlConfig = errors.New("invalid mysql Config")

type Mysql struct {
	Username     string `mapstructure:"username" json:"username" yaml:"username"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	Host         string `mapstructure:"host" json:"host" yaml:"host"`
	Dbname       string `mapstructure:"db" json:"db" yaml:"db"`
	Config       string `mapstructure:"config" json:"config" yaml:"config"`
	MaxIdleConns int    `mapstructure:"max-idle-conns" json:"maxIdleConns" yaml:"max-idle-conns"`
	MaxOpenConns int    `mapstructure:"max-open-conns" json:"maxOpenConns" yaml:"max-open-conns"`
	LogMode      bool   `mapstructure:"log" json:"log" yaml:"log"`
}

func (p *Mysql) isValid() bool {
	if p.Host == "" || p.Username == "" || p.Password == "" ||
		p.Dbname == "" {
		return false
	}

	return true
}

// 初始化数据库并产生数据库全局变量
func initMysql() {
	admin := sysConf().Mysql
	if !admin.isValid() {
		panic(InvalidMysqlConfig)
	}
	if db, err := gorm.Open("mysql", admin.Username+":"+admin.Password+"@("+admin.Host+")/"+admin.Dbname+"?"+admin.Config); err != nil {
		panic(fmt.Sprintf("MySQL链接异常, %v", err))
	} else {
		_db = db
		DB().DB().SetMaxIdleConns(admin.MaxIdleConns)
		DB().DB().SetMaxOpenConns(admin.MaxOpenConns)
		DB().LogMode(admin.LogMode)
	}
}
