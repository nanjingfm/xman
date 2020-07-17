package xman

import (
	"github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = ginkgo.Describe("initSysConfig", func() {
	originalFile := _sysConfigFile
	ginkgo.BeforeEach(func() {
		_sysConfigFile = "./testdata/config/sys.yaml"
	})
	ginkgo.Context("init sys config", func() {
		ginkgo.It("return mysql username", func() {
			initSysConfig()
			Expect(_config.Mysql.Username).To(Equal("root"))
		})
	})
	ginkgo.AfterEach(func() {
		_sysConfigFile = originalFile
	})
})

var _ = ginkgo.Describe("initAppConfig", func() {
	appFile := _appConfigFile
	ginkgo.BeforeEach(func() {
		_appConfigFile = "./testdata/config/app.yaml"
	})
	ginkgo.Context("init app config", func() {
		ginkgo.It("return switch.mod", func() {
			initAppConfig()
			Expect(_appConfig.Get("switch.mod")).To(Equal("prod"))
		})
	})
	ginkgo.AfterEach(func() {
		_appConfigFile = appFile
	})
})
