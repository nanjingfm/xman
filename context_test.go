package xman

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestXman(t *testing.T) {
	RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Xman Suite")
}

var _ = ginkgo.Describe("Context", func() {
	ginkgo.Describe("empty locale context", func() {
		ginkgo.It("return default locale", func() {
			ginCtx := &gin.Context{}
			c := ParseContext(ginCtx)
			Expect(c.Locale).To(Equal(1001))
		})
	})
})
