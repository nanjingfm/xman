package xman

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/unknwon/i18n"
	"testing"
)

func TestECode_Error(t *testing.T) {
	tests := []struct {
		name string
		Code int
		Msg  string
		want string
	}{
		{
			name: "with err msg",
			Code: 10009,
			Msg:  "err msg",
			want: "err msg",
		},
		{
			name: "without err msg",
			Code: 10009,
			want: "code.10009",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ECode{
				Code: tt.Code,
				Msg:  tt.Msg,
			}
			if got := p.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestECode_GetCode(t *testing.T) {
	assert.Equal(t, -2, ECodeUnknownErr.GetCode())
	assert.Equal(t, -1, ECodeSystemErr.GetCode())
	assert.Equal(t, 1001, ECodeSuccess.GetCode())
}

func TestECode_GetMsg(t *testing.T) {
	type fields struct {
		Code int
		Msg  string
	}
	type args struct {
		ext []string
	}
	tests := []struct {
		name string
		Code int
		Msg  string
		args args
		want string
	}{
		{
			name: "with err msg",
			Code: 10009,
			Msg:  "err msg",
			want: "err msg",
		},
		{
			name: "without err msg",
			Code: 10009,
			want: "code.10009",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := ECode{
				Code: tt.Code,
				Msg:  tt.Msg,
			}
			if got := p.GetMsg(tt.args.ext...); got != tt.want {
				t.Errorf("GetMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCode(t *testing.T) {
	eCode := NewCode(123)
	assert.Equal(t, 123, eCode.GetCode())
}

func TestNewErrorCode(t *testing.T) {
	code := NewErrorCode(nil)
	assert.Equal(t, ECodeSuccess, code)

	code = NewErrorCode(errors.New("xxx"))
	assert.Equal(t, ECodeSystemErr.GetCode(), code.GetCode())
	assert.Equal(t, "xxx", code.GetMsg())

	code = NewErrorCode(NewCode(666))
	assert.Equal(t, 666, code.GetCode())
	assert.Equal(t, "code.666", code.GetMsg())

	err := errors.New("original error")
	w := fmt.Errorf("Wrap error:%w", err)
	code = NewErrorCode(w)
	assert.Equal(t, ECodeSystemErr.GetCode(), code.GetCode())
	assert.Equal(t, "Wrap error:original error", code.GetMsg())
}

func TestECode_getCodeKey(t *testing.T) {
	c1 := NewCode(111)
	assert.Equal(t, "code.111", c1.getCodeKey())

	config := I18nOptions{
		Format: "%s.ini",
		Files: map[string][]byte{
			"zh-CN": []byte("aaa=111"),
		},
		Directory:   "./testdata/config/locale",
		DefaultLang: LangZhCN,
		Langs:       []string{LangZhCN, LangZhTW},
		Names:       []string{"简体中文", "繁体中文"},
	}

	SetLocale(config)
	c2 := NewCode(666)
	s := c2.GetLocaleMsg(Locale{i18n.Locale{Lang: LangZhCN}})
	assert.Equal(t, "成功", s)
}
