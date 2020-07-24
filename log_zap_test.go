package xman

import (
	"go.uber.org/zap/zapcore"
	"testing"
)

func Test_initLogZap(t *testing.T) {
	_config.System.Env = EnvDev
	_config.Log = LogConfig{
		LogFile:  "./test.log",
		LogLevel: zapcore.DebugLevel,
	}
	initLogZap()

	LogError("test-error", "key1", "value1", "key2", 2)
	LogInfo("test-info", "key1", "value1", "key2", 2, "")
}
