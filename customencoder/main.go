package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SysTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Jan  2 15:04:05"))
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func main() {
	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding: "console",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "message",
			LevelKey:    "level",
			TimeKey:     "time",
			EncodeLevel: zapcore.CapitalColorLevelEncoder,
			EncodeTime:  zapcore.ISO8601TimeEncoder,
		},
		OutputPaths: []string{"stderr"},
	}

	fmt.Printf("\n*** Using standard ISO8601 time encoder\n\n")

	logger, _ := cfg.Build()
	logger.Info("This should have an ISO8601 based time stamp")

	fmt.Printf("\n*** Using a custom time encoder\n\n")

	cfg.EncoderConfig.EncodeTime = SysTimeEncoder

	logger, _ = cfg.Build()
	logger.Info("This should have a syslog style time stamp")

	fmt.Printf("\n*** Using a custom level encoder\n\n")

	cfg.EncoderConfig.EncodeLevel = CustomLevelEncoder

	logger, _ = cfg.Build()
	logger.Info("This should have a interesting level name")
}
