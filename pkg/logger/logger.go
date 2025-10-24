package logger

import (
	"os"

	"github.com/In-the-name-and-glory-of-God/entrepreneur-pastoral/pkg/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(env string) *zap.SugaredLogger {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	logLevel := zapcore.InfoLevel
	if env == constants.DEVELOPMENT {
		logLevel = zapcore.DebugLevel
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		logLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).Sugar()
}
