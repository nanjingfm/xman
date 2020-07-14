package xman

func InitServer() {
	initConfig()
	switch sysConf().System.DbType {
	case "mysql":
		initMysql()
	case "sqlite":
		initSqlite() // sqlite需要gcc支持 windows用户需要自行安装gcc 如需使用打开注释即可
	default:
		initMysql()
	}

	if sysConf().System.UseMultipoint {
		// 初始化redis服务
		initRedis()
	}

	initI18n() // 初始化多语言
}
