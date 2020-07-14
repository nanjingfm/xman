package xman

import (
	"bytes"
	"net/http"
	"path"
	"strings"
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
	CaptchaId string `json:"captchaId"`
	PicPath   string `json:"picPath"`
}

// CaptchaVerify 验证验证码有效性
func CaptchaVerify(id string, digits string) bool {
	return captcha.VerifyString(id, digits)
}

// registerCaptchaRouter 注册路由
func registerCaptchaRouter(Router *gin.Engine) (R gin.IRoutes) {
	BaseRouter := Router.Group("base")
	{
		BaseRouter.POST("captcha", captchaHandle)
		BaseRouter.GET("captcha/:captchaId", captchaImg)
	}
	return BaseRouter
}

func captchaHandle(c *gin.Context) {
	captchaId := captcha.NewLen(sysConf().Captcha.KeyLong)
	Return(c, ECodeSuccess, SysCaptchaResponse{
		CaptchaId: captchaId,
		PicPath:   "/base/captcha/" + captchaId + ".png",
	}, "验证码获取成功")
}

func captchaImg(c *gin.Context) {
	ginCaptchaServeHTTP(c.Writer, c.Request)
}

// 这里需要自行实现captcha 的gin模式
func ginCaptchaServeHTTP(w http.ResponseWriter, r *http.Request) {
	dir, file := path.Split(r.URL.Path)
	ext := path.Ext(file)
	id := file[:len(file)-len(ext)]
	if ext == "" || id == "" {
		http.NotFound(w, r)
		return
	}
	if r.FormValue("reload") != "" {
		captcha.Reload(id)
	}
	lang := strings.ToLower(r.FormValue("lang"))
	download := path.Base(dir) == "download"
	if serve(w, r, id, ext, lang, download, sysConf().Captcha.ImgWidth, sysConf().Captcha.ImgHeight) == captcha.ErrNotFound {
		http.NotFound(w, r)
	}
}

func serve(w http.ResponseWriter, r *http.Request, id, ext, lang string, download bool, width, height int) error {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	var content bytes.Buffer
	switch ext {
	case ".png":
		w.Header().Set("Content-Type", "image/png")
		_ = captcha.WriteImage(&content, id, width, height)
	case ".wav":
		w.Header().Set("Content-Type", "audio/x-wav")
		_ = captcha.WriteAudio(&content, id, lang)
	default:
		return captcha.ErrNotFound
	}

	if download {
		w.Header().Set("Content-Type", "application/octet-stream")
	}
	http.ServeContent(w, r, id+ext, time.Time{}, bytes.NewReader(content.Bytes()))
	return nil
}
