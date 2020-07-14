package xman

import (
	"github.com/gin-gonic/gin"
)

// Context 上下文
type Context struct {
	Locale // 语言包
	//Data   map[string]interface{} // 自定义数据
}

func ParseContext(c *gin.Context) *Context {
	ctx := &Context{}
	ctx.Locale = _defaultLocale
	if data, exist := c.Get(_localeContextKey); exist {
		if l, ok := data.(Locale); ok {
			ctx.Locale = l
		}
	}

	return ctx
}
