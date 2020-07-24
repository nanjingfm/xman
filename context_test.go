package xman

import (
	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/unknwon/i18n"
)

var _ = ginkgo.Describe("Context", func() {
	ginkgo.Describe("empty locale context", func() {
		ginkgo.It("return default locale", func() {
			ginCtx := &gin.Context{}
			c := ParseContext(ginCtx)
			Expect(c.Locale).To(Equal(_globalLocale))
		})

		ginkgo.It("return ch locale", func() {
			ginCtx := &gin.Context{}
			ginCtx.Set(_localeContextKey, Locale{Locale: i18n.Locale{Lang: LangZhCN}})
			c := ParseContext(ginCtx)
			Expect(c.Locale.Lang).To(Equal(LangZhCN))
		})
	})
})
