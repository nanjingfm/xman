package xman

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func Return(c *gin.Context, code ECode, data interface{}, msg ...string) {
	if code.GetCode() == 0 {
		code = ECodeUnknownErr
	}
	codeStr := code.GetLocaleMsg(ParseContext(c).Locale)
	c.JSON(http.StatusOK, Response{
		code.GetCode(),
		data,
		codeStr,
	})
}
