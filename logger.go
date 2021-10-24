package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var logger *zap.SugaredLogger

func InitLogger(logFilePath string) {
	writerSyncer := getLogWriter(logFilePath)
	encoder := getEncoder(logFilePath)

	core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

	zapLogger := zap.New(core)
	logger = zapLogger.Sugar()

	zap.ReplaceGlobals(zapLogger)
}

func getEncoder(logFilePath string) zapcore.Encoder {
	//if no log file path then output should be in console
	if logFilePath == "" {
		return zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	}
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter(logFilePath string) zapcore.WriteSyncer {
	var file *os.File
	if logFilePath != "" {
		file, _ = os.Create(logFilePath)
	} else {
		file = os.Stdout
	}
	return zapcore.AddSync(file)
}
