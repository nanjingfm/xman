package xman

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unknwon/i18n"
	"testing"
)

func TestMain(t *testing.M) {
	_ = i18n.SetMessage(LangZhCN, "./testdata/config/locale/zh-CN.ini")
	_defaultOptions = I18nOptions{
		Format:      "%s.ini",
		Directory:   "./testdata/config/locale",
		Langs:       []string{LangZhCN},
		DefaultLang: LangZhCN,
		Names:       []string{"简体中文"},
		Parameter:   "lang",
	}
	t.Run()
}

func TestXman(t *testing.T) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Xman Suite")
}
