package xman

// sqlite需要gcc支持 windows用户需要自行安装gcc 如需使用打开注释即可

// 感谢 sqlitet提供者 [rikugun] 作者github： https://github.com/rikugun

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Sqlite struct {
	Username string `mapstructure:"username" json:"username" yaml:"username"`
	Password string `mapstructure:"password" json:"password" yaml:"password"`
	Path     string `mapstructure:"path" json:"path" yaml:"path"`
	Config   string `mapstructure:"config" json:"config" yaml:"config"`
	LogMode  bool   `mapstructure:"log-mode" json:"logMode" yaml:"log-mode"`
}

// 初始化数据库并产生数据库全局变量
func initSqlite() {
	admin := sysConf().Sqlite
	if db, err := gorm.Open("sqlite3", fmt.Sprintf("%s?%s", admin.Path, admin.Config)); err != nil {
		LogError("DEFAULTDB数据库启动异常", err)
	} else {
		_db = db
		DB().LogMode(admin.LogMode)
	}
}
