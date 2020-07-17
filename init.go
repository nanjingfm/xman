package xman

func InitServer() {
	initConfig()
	switch sysConf().System.DbType {
	case "mysql":
		initMysql()
	default:
		initMysql()
	}

	if sysConf().System.UseMultipoint {
		// 初始化redis服务
		initRedis()
	}

	initI18n() // 初始化多语言
}
