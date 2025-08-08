package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder

	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
	logFile, err := os.OpenFile("./log/app.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	fileWriter := zapcore.AddSync(logFile)

	fileLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	fileCore := zapcore.NewCore(fileEncoder, fileWriter, fileLevel)

	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleCore := zapcore.NewCore(consoleEncoder, os.Stdout, zapcore.DebugLevel)

	core := zapcore.NewTee(fileCore, consoleCore)

	logger := zap.New(core, zap.AddCaller())
	return logger
}
