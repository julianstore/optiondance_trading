package log

import "go.uber.org/zap"

func Error(msg string, fields ...zap.Field) {
	zap.L().Error(msg, fields...)
}

func ErrorInfo(msg string, err error) {
	zap.L().Info(msg, zap.Error(err))
}
