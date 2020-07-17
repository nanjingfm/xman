package xman

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

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

// 初始化数据库并产生数据库全局变量
func initMysql() {
	admin := sysConf().Mysql
	if db, err := gorm.Open("mysql", admin.Username+":"+admin.Password+"@("+admin.Host+")/"+admin.Dbname+"?"+admin.Config); err != nil {
		panic(fmt.Sprintf("MySQL启动异常, %v", err))
	} else {
		_db = db
		DB().DB().SetMaxIdleConns(admin.MaxIdleConns)
		DB().DB().SetMaxOpenConns(admin.MaxOpenConns)
		DB().LogMode(admin.LogMode)
	}
}
