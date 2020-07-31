package xman

import (
	"strconv"
)

var (
	ECodeSuccess    = NewCode(1001)
	ECodeSystemErr  = NewCode(-1) // 系统错误
	ECodeUnknownErr = NewCode(-2) // 未知错误
	ECodeCaptchaErr = NewCode(-3) // 验证码错误
	ECodeParamErr   = NewCode(-4) // 参数错误
)

type Coder interface {
}

func NewCode(code int) ECode {
	return ECode{Code: code}
}

func NewCodeMsg(code int, msg string) ECode {
	return ECode{Code: code, Msg: msg}
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
	return p.Msg
}

func (p ECode) getCodeKey() string {
	codeStr := strconv.Itoa(p.Code)
	return "code." + codeStr
}

func (p ECode) GetLocaleMsg(l Locale) string {
	key := p.getCodeKey()
	str := l.Tr(key)
	if str == key {
		if msg := p.GetMsg(); msg != "" {
			return msg
		} else {
			return key
		}
	}

	return str
}
