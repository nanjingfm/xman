package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nanjingfm/xman"
)

func main() {
	e := gin.Default()
	xman.RegisterCaptchaRouter(e)
	e.GET("/check-captcha", xman.CaptchaAuth(), func(context *gin.Context) {
		xman.Return(context, xman.ECodeSuccess, nil)
	})
	e.Run(":9998")
}
