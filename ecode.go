package xman

import (
	"strconv"
)

var (
	ECodeSuccess    = NewCode(0)
	ECodeSystemErr  = NewCode(7)
	ECodeUnknownErr = NewCode(10) // 未知错误
)

type Coder interface {
}

func NewCode(code int) *ECode {
	return &ECode{Code: code}
}

func NewErrorCode(err error) *ECode {
	return &ECode{Code: 1003, Msg: err.Error()}
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
	codeMsg := _defaultLocale.Tr("code." + codeStr)
	if codeMsg != "" {
		return codeMsg
	}
	return "code: " + codeStr
}
