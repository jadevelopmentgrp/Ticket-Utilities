package observability

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Configure(json bool, logLevel zapcore.Level) (*zap.Logger, error) {
	if json {
		loggerConfig := zap.NewProductionConfig()
		loggerConfig.Level.SetLevel(logLevel)

		return loggerConfig.Build(
			zap.AddCaller(),
			zap.AddStacktrace(zap.ErrorLevel),
			zap.WrapCore(ZapAdapter()),
		)
	} else {
		loggerConfig := zap.NewDevelopmentConfig()
		loggerConfig.Level.SetLevel(logLevel)
		loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		return loggerConfig.Build(zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	}
}
