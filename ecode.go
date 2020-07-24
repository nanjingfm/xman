package xman

import (
	"strconv"
)

var (
	ECodeSuccess    = NewCode(1001)
	ECodeSystemErr  = NewCode(-1) // 系统错误
	ECodeUnknownErr = NewCode(-2) // 未知错误
	ECodeCaptchaErr = NewCode(-3) // 验证码错误
)

type Coder interface {
}

func NewCode(code int) ECode {
	return ECode{Code: code}
}

func NewErrorCode(err error) ECode {
	if err == nil {
		return ECodeSuccess
	}
	var (
		eCodeErr ECode
		ok       bool
	)
	if eCodeErr, ok = err.(ECode); ok {
		return eCodeErr
	}
	eCodeErr = ECodeSystemErr
	eCodeErr.Msg = err.Error()
	return eCodeErr
}

type ECode struct {
	Code int
	Msg  string
}

func (p ECode) GetCode() int {
	return p.Code
}

func (p ECode) GetMsg(ext ...string) string {
	return p.Error()
}

func (p ECode) Error() string {
	if p.Msg != "" {
		return p.Msg
	}

	codeStr := strconv.Itoa(p.Code)
	return "code." + codeStr
}
