package xman

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unknwon/i18n"
	"testing"
)

func TestMain(t *testing.M) {
	_ = i18n.SetMessage(LangZhCN, "./testdata/locale/locale_zh-CN.ini")
	t.Run()
}

func TestXman(t *testing.T) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Xman Suite")
}
