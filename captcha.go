package xman

import (
	"bytes"
	"net/http"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

type Captcha struct {
	KeyLong   int `mapstructure:"key-long" json:"keyLong" yaml:"key-long"`
	ImgWidth  int `mapstructure:"img-width" json:"imgWidth" yaml:"img-width"`
	ImgHeight int `mapstructure:"img-height" json:"imgHeight" yaml:"img-height"`
}

type SysCaptchaResponse struct {
	CaptchaId string `json:"captcha_id"`
	PicPath   string `json:"pic_path"`
}

// CaptchaVerify 验证验证码有效性
func CaptchaVerify(id string, digits string) bool {
	return captcha.VerifyString(id, digits)
}

// registerCaptchaRouter 注册路由
func registerCaptchaRouter(Router *gin.Engine) (R gin.IRoutes) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.GET("captcha", captchaHandle)
		BaseRouter.GET("captcha-img", captchaImg)
	}
	return BaseRouter
}

func captchaHandle(c *gin.Context) {
	captchaId := captcha.NewLen(sysConf().Captcha.KeyLong)
	Return(c, ECodeSuccess, SysCaptchaResponse{
		CaptchaId: captchaId,
	}, "验证码获取成功")
}

func captchaImg(c *gin.Context) {
	var (
		captchaId string
		reload bool
	)
	captchaId = c.GetString("captchaId")
	if captchaId == "" {
		c.Writer.WriteHeader(http.StatusNotFound)
		return
	}
	reload = c.GetBool("reload")
	if reload {
		captcha.Reload(captchaId)
	}
	width := sysConf().Captcha.ImgWidth
	height := sysConf().Captcha.ImgHeight
	header := c.Writer.Header()
	header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	header.Set("Pragma", "no-cache")
	header.Set("Expires", "0")
	header.Set("Content-Type", "image/png")
	header.Set("Content-Type", "application/octet-stream")
	var content bytes.Buffer
	_ = captcha.WriteImage(&content, captchaId, width, height)
	http.ServeContent(c.Writer, c.Request, captchaId+".ext", time.Time{}, bytes.NewReader(content.Bytes()))
}
