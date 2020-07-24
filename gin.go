package xman

import (
	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"net/http"
)

type Engine struct {
	*gin.Engine
}

func NewEngine() *Engine {
	var r = gin.Default()
	// Router.Use(middleware.LoadTls())  // 打开就能玩https了

	r.Use(Cors()) // 跨域
	r.Use(I18n()) // 多语言

	return &Engine{Engine: r}
}

// AddRoute 注册路由
func (e *Engine) AddRoute(r func(r *gin.Engine)) {
	r(e.Engine)
}

// Run 启动
func (e *Engine) Run(addr ...string) {
	LogDebug("server run success on ", addr)
	err := e.Engine.Run(addr...)
	LogError("Run http", "err", err)
}

func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "localhost:443",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			// 如果出现错误，请不要继续
			LogError("load Tls error")
			return
		}
		// 继续往下处理
		c.Next()
	}
}

// 处理跨域请求,支持options访问
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		// 处理请求
		c.Next()
	}
}
