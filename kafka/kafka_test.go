package kafka

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var logger *zap.Logger

func init() {
	var err error
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = TimeEncoder
	logger, err = cfg.Build()
	if err != nil {
		panic(err)
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format("2006-01-02 15:04:05") + "]")
}
