package xman

import (
	"github.com/unknwon/i18n"
	"go.uber.org/zap"
)

var (
	_globalLocale = Locale{Locale: i18n.Locale{Lang: LangZhCN}}
	_globalLogger = NewLogDev(_defaultLogConfig)
)

func SetLogger(logger *zap.SugaredLogger) {
	if logger == nil {
		panic("SugaredLogger nil error")
	}
	_globalLogger = logger
}

func SetLocale(config I18nOptions) {
	_defaultOptions = config
	_globalLocale = newI18n(config)
}

func Tr(format string, args ...interface{}) string {
	return _globalLocale.Tr(format, args)
}

func LogError(msg string, args ...interface{}) {
	_globalLogger.Errorw(msg, args...)
}

func LogInfo(msg string, args ...interface{}) {
	_globalLogger.Infow(msg, args...)
}

func LogDebug(msg string, args ...interface{}) {
	_globalLogger.Debugw(msg, args...)
}
