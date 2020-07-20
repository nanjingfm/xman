package xman

import (
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"github.com/tidwall/gjson"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturn(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	Return(c, ECodeSuccess, nil)

	body := w.Body.String()
	assert.Equal(t, http.StatusOK, c.Writer.Status())
	assert.Equal(t, int64(1001), gjson.Get(body, "code").Int())

	w = httptest.NewRecorder()
	c, _ = gin.CreateTestContext(w)
	Return(c, NewCode(0), nil)

	body = w.Body.String()
	assert.Equal(t, http.StatusOK, c.Writer.Status())
	assert.Equal(t, int64(ECodeUnknownErr.GetCode()), gjson.Get(body, "code").Int())
}
