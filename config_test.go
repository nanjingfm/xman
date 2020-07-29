package xman

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const _testSysConfigFile = "./testdata/config/sys.yaml"
const _testAppConfigFile = "./testdata/config/app.yaml"

var _ = ginkgo.Describe("initSysConfig", func() {
	ginkgo.Context("init sys config", func() {
		ginkgo.It("return mysql username", func() {
			sysConfig := NewSysConfig(_testSysConfigFile)
			Expect(sysConfig.DB.DbType).To(Equal("mysql"))
		})
	})
})

var _ = ginkgo.Describe("initAppConfig", func() {
	ginkgo.Context("init app config", func() {
		ginkgo.It("return switch.mod", func() {
			appConfig := NewAppConfig(_testAppConfigFile)
			Expect(appConfig.Get("switch.mod")).To(Equal("prod"))
		})
	})
})
