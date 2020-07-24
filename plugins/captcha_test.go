package plugins

import (
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestFetchCaptcha(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	RegisterCaptchaRouter(engine)
	r := httptest.NewRequest("GET", "/base/captcha", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Header().Get("Content-type"), "image/png"))
	assert.True(t, w.Header().Get(_captchaID) != "")
}

func TestCaptchaVerify(t *testing.T) {
	assert.False(t, CaptchaVerify("xxxxxx", "123"))

	s := captcha.NewMemoryStore(10, time.Minute)
	captcha.SetCustomStore(s)
	captchaID := captcha.New()
	digits := s.Get(captchaID, false)
	ns := make([]byte, 0)
	for i := range digits {
		d := digits[i]
		switch {
		case '0' <= d+'0' && d+'0' <= '9':
			ns = append(ns, d+'0')
		case d == ' ' || d == ',':
			// ignore
		default:
			ns = append(ns, d)
		}
	}
	assert.True(t, CaptchaVerify(captchaID, string(ns)))
}

func TestCaptchaAuth(t *testing.T) {
	s := captcha.NewMemoryStore(10, time.Minute)
	captcha.SetCustomStore(s)
	captchaID := captcha.New()
	digits := s.Get(captchaID, false)
	ns := make([]byte, 0)
	for i := range digits {
		d := digits[i]
		switch {
		case '0' <= d+'0' && d+'0' <= '9':
			ns = append(ns, d+'0')
		case d == ' ' || d == ',':
			// ignore
		default:
			ns = append(ns, d)
		}
	}
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request, _ = http.NewRequest("GET", "/", nil)
	c.Request.Header.Set(_captchaID, captchaID)
	c.Request.Header.Set(_captchaValue, string(ns))
	f := CaptchaAuth()
	f(c)
	assert.Equal(t, 200, c.Writer.Status())
}
