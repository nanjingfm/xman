package xman

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func Test_isFile(t *testing.T) {
	tempPath := os.TempDir()
	assert.False(t, isFile(tempPath))

	f, err := ioutil.TempFile(tempPath, "xxx.log")
	defer os.Remove(f.Name())
	assert.Nil(t, err)
	assert.True(t, isFile(f.Name()))
}

func Test_initLocales(t *testing.T) {
	matcher := initLocales(_defaultOptions)
	assert.NotNil(t, matcher)

	initLocales(I18nOptions{
		Format:    "%s.ini",
		Directory: "./testdata/config/locale/",
		Files: map[string][]byte{
			"zh-CN": []byte("aaa=111"),
		},
		Langs: []string{LangZhCN, LangZhTW},
		Names: []string{"简体中文", "繁体中文"},
	})
}

func TestLocale_Language(t *testing.T) {
	l := Locale{}
	l.Lang = LangZhCN
	assert.Equal(t, LangZhCN, l.Language())
}

func Test_initI18n(t *testing.T) {
	_config.I18n = I18nOptions{
		Format:    "%s.ini",
		Directory: "./testdata/config/locale/",
		Files: map[string][]byte{
			"zh-CN": []byte("aaa=111"),
		},
		Langs: []string{LangZhCN, LangZhTW},
		Names: []string{"简体中文", "繁体中文"},
	}
	initI18n()
	assert.Equal(t, LangZhCN, _defaultLocale.Language())
	assert.NotNil(t, _defaultOptions)
}

func TestI18n(t *testing.T) {
	_config.I18n = I18nOptions{
		Format: "%s.ini",
		Files: map[string][]byte{
			"zh-CN": []byte("aaa=111"),
		},
		Directory:   "./testdata/config/locale",
		DefaultLang: LangZhCN,
		Langs:       []string{LangZhCN, LangZhTW},
		Names:       []string{"简体中文", "繁体中文"},
	}
	initI18n()

	gin.SetMode(gin.ReleaseMode)
	e := gin.New()
	e.Use(I18n())
	e.GET("/test", func(context *gin.Context) {
		Return(context, ECodeSuccess, nil)
	})

	// 测试默认语言
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/test", nil)
	e.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "成功", gjson.Get(w.Body.String(), "msg").String())

	// 测试自定义语言
	w2 := httptest.NewRecorder()
	r2 := httptest.NewRequest("GET", "/test", nil)
	r2.Header.Add("Cookie", "lang=zh-TW")
	e.ServeHTTP(w2, r2)
	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "胜利", gjson.Get(w2.Body.String(), "msg").String())

	// 测试不存在的语言
	w3 := httptest.NewRecorder()
	r3 := httptest.NewRequest("GET", "/test", nil)
	r3.Header.Add("Cookie", "lang=zh-TWxx")
	_defaultOptions.Redirect = true
	e.ServeHTTP(w3, r3)
	assert.Equal(t, http.StatusOK, w3.Code)
	assert.True(t, strings.Contains(w3.Header().Get("Set-Cookie"), LangZhCN))
}
