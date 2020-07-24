package xman

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	LogFile  string        `json:"log-file" yaml:"log-file"`
	LogLevel zapcore.Level `json:"log-level" yaml:"log-level"`
}

func initLogZap() {
	var (
		err    error
		logger *zap.Logger
	)
	logConfig := sysConf().Log
	config := &zap.Config{}
	if IsDev() {
		config.Development = true
		config.Encoding = "console"
		config.OutputPaths = []string{"stderr", logConfig.LogFile}
		config.ErrorOutputPaths = []string{"stderr"}
	} else {
		config.Development = false
		config.Encoding = "json"
		config.OutputPaths = []string{"stderr"}
		config.ErrorOutputPaths = []string{"stderr"}
	}

	config.Level = zap.NewAtomicLevelAt(logConfig.LogLevel)
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	logger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("初始化logger失败，err:%v", err))
	}
	_logger = logger.Sugar()
}
