package xman

import (
	"go.uber.org/zap/zapcore"
	"testing"
)

func Test_initLogZap(t *testing.T) {
	config := LogConfig{
		LogFile:  "./test.log",
		LogLevel: zapcore.DebugLevel,
	}
	NewLogDev(config)

	LogError("test-error", "key1", "value1", "key2", 2)
	LogInfo("test-info", "key1", "value1", "key2", 2)

	config2 := LogConfig{
		LogFile:  "./test.log",
		LogLevel: zapcore.DebugLevel,
	}
	NewLogProd(config2)

	LogError("test-error", "key1", "value1", "key2", 2)
	LogInfo("test-info", "key1", "value1", "key2", 2)
}
