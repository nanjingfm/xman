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

func ReturnParamErr(c *gin.Context)  {
	Return(c, ECodeParamErr, nil)
}

func Return(c *gin.Context, err error, data interface{}, msg ...string) {
	code, ok := err.(ECode)
	if !ok {
		code = NewErrorCode(err)
	}
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
