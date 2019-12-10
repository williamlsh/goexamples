package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	//cfg := zap.Config{
	//	Level: zap.NewAtomicLevelAt(zapcore.InfoLevel),
	//	Encoding: "json",
	//}
	//logger, err := cfg.Build()
	//if err != nil {
	//	panic(err)
	//}

	fmt.Printf("\n*** Using a JSON encoder, at debug level, sending output to stdout, no key specified\n\n")
	logger, _ := zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
	}.Build()

	logger.Debug("This is a debug message")
	logger.Info("This is an info message")
	logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))

	fmt.Printf("\n*** Using a JSON encoder, at debug level, sending output to stdout, message key only specified\n\n")
	logger, _ = zap.Config{
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding:    "json",
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",
		},
	}.Build()

	logger.Info("This is a debug message")
	logger.Info("This is an info message")
	logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))

	fmt.Printf("\n*** Using a JSON encoder, at debug level, sending output to stdout, all possible keys specified\n\n")
	cfg := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "caller",
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	logger, _ = cfg.Build()
	logger.Info("This is a debug message")
	logger.Info("This is an info message")
	logger.Info("This is a info message with fields", zap.String("region", "us-west"), zap.Int("id", 7))

	fmt.Printf("\n*** Same logger with console logging enabled instead\n\n")
	logger.WithOptions(zap.WrapCore(func(zapcore.Core) zapcore.Core {
		return zapcore.NewCore(zapcore.NewConsoleEncoder(cfg.EncoderConfig), zapcore.AddSync(os.Stderr), zapcore.DebugLevel)
	})).Info("This is an info message")
}
