package logger

import (
	"os"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(app config.Application) *zap.SugaredLogger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	logLevel := zapcore.InfoLevel
	if app.IsDevelopment() {
		logLevel = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		logLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}
