package xman

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var _defaultLogConfig = LogConfig{LogLevel: zap.InfoLevel}

type LogConfig struct {
	LogFile  string        `mapstructure:"log-file" yaml:"log-file"`
	LogLevel zapcore.Level `mapstructure:"log-level" yaml:"log-level"`
}

var DevLogConfig = zap.Config{
	Development:      true,
	Encoding:         "console",
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
}

var ProdLogConfig = zap.Config{
	Development:      false,
	Encoding:         "json",
	OutputPaths:      []string{"stderr"},
	ErrorOutputPaths: []string{"stderr"},
}

func NewLogDev(logConfig LogConfig) *zap.SugaredLogger {
	var (
		err    error
		logger *zap.Logger
	)

	config := DevLogConfig
	if logConfig.LogFile != "" {
		config.OutputPaths = append(config.OutputPaths, logConfig.LogFile)
	}
	config.Level = zap.NewAtomicLevelAt(logConfig.LogLevel)
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	logger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("初始化logger失败，err:%v", err))
	}
	return logger.Sugar()
}

func NewLogProd(logConfig LogConfig) *zap.SugaredLogger {
	var (
		err    error
		logger *zap.Logger
	)

	config := ProdLogConfig
	config.OutputPaths = append(config.OutputPaths, logConfig.LogFile)
	config.Level = zap.NewAtomicLevelAt(logConfig.LogLevel)
	config.EncoderConfig = zap.NewProductionEncoderConfig()
	logger, err = config.Build()
	if err != nil {
		panic(fmt.Sprintf("初始化logger失败，err:%v", err))
	}
	return logger.Sugar()
}
