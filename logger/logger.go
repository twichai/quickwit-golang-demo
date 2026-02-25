package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New() (*zap.Logger, error) {
	file, err := os.OpenFile("./logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewProductionConfig()

	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	cfg.EncoderConfig.TimeKey = "creationDate"

	cfg.OutputPaths = []string{}
	cfg.ErrorOutputPaths = []string{}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg.EncoderConfig),
		zapcore.AddSync(file),
		zap.InfoLevel,
	)

	return zap.New(core), nil
}
