package plugins

import (
	"bytes"
	"github.com/nanjingfm/xman"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

const (
	_captchaID     = "x-captcha-id"
	_captchaValue  = "x-captcha-value"
	_defaultWidth  = 300
	_defaultHeight = 100
)

type Captcha struct {
	KeyLong int `yaml:"key-long"`
	CaptchaSize
}

type CaptchaSize struct {
	Width  int `form:"width" yaml:"width"`
	Height int `form:"height" yaml:"height"`
}

type SysCaptchaResponse struct {
	CaptchaId string `json:"captcha_id"`
	PicPath   string `json:"pic_path"`
}

// CaptchaVerify 验证验证码有效性
func CaptchaVerify(id string, digits string) bool {
	return captcha.VerifyString(id, digits)
}

// RegisterCaptchaRouter 注册路由
func RegisterCaptchaRouter(Router *gin.Engine) (R gin.IRoutes) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", captchaHandle)
	}
	return BaseRouter
}

// captchaHandle 生成验证码
func captchaHandle(c *gin.Context) {
	size := CaptchaSize{}
	_ = c.BindQuery(&size)
	if size.Width == 0 {
		size.Width = _defaultWidth
	}

	if size.Height == 0 {
		size.Height = _defaultHeight
	}

	if size.Width == 0 {
		size.Width = 200
	}

	if size.Height == 0 {
		size.Height = 50
	}

	captchaId := captcha.NewLen(4) // TODO
	header := c.Writer.Header()
	header.Set(_captchaID, captchaId)

	header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	header.Set("Pragma", "no-cache")
	header.Set("Expires", "0")
	header.Set("Content-Type", "image/png")

	var content bytes.Buffer
	_ = captcha.WriteImage(&content, captchaId, size.Width, size.Height)
	http.ServeContent(c.Writer, c.Request, captchaId+".png", time.Time{}, bytes.NewReader(content.Bytes()))
}

// CaptchaAuth 验证图片验证码
func CaptchaAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		captchaID := c.GetHeader(_captchaID)
		digits := c.GetHeader(_captchaValue)
		c.Writer.Header().Del(_captchaValue)
		if captchaID == "" || digits == "" || !CaptchaVerify(captchaID, digits) {
			xman.Return(c, xman.ECodeCaptchaErr, nil)
			c.Abort()
			return
		}
		// 删除header
		c.Writer.Header().Del(_captchaID)
		// 处理请求
		c.Next()
	}
}
