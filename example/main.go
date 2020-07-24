package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nanjingfm/xman"
	"github.com/nanjingfm/xman/plugins"
)

func main() {
	e := gin.Default()
	plugins.RegisterCaptchaRouter(e)
	e.GET("/check-captcha", plugins.CaptchaAuth(), func(context *gin.Context) {
		xman.Return(context, xman.ECodeSuccess, nil)
	})
	e.Run(":9998")
}
