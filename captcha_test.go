package xman

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchCaptcha(t *testing.T) {
	engine := gin.New()
	registerCaptchaRouter(engine)
	r := httptest.NewRequest("GET", "/base/captcha", nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Header().Get("Content-type"),"application/json"))
	data := w.Body.String()

	assert.Equal(t, int64(1001), gjson.Get(data, "code").Int())
	assert.NotEmpty(t, gjson.Get(data, "data.captcha_id").String())
}

func TestCaptchaVerify(t *testing.T) {
	engine := gin.New()
	registerCaptchaRouter(engine)
	id := captcha.New()
	r := httptest.NewRequest("GET", "/base/captcha/"+id, nil)
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, r)
	data := w.Body.String()
	spew.Dump(data)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.True(t, strings.Contains(w.Header().Get("Content-type"),"application/json"))
}
