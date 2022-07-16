package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	loggerSafeGuard sync.Once
	logger          *zap.SugaredLogger
)

func InitLogger(logFilePath string) {
	loggerSafeGuard.Do(func() {
		writerSyncer := getLogWriter(logFilePath)
		encoder := getEncoder(logFilePath)

		core := zapcore.NewCore(encoder, writerSyncer, zapcore.DebugLevel)

		zapLogger := zap.New(core)
		logger = zapLogger.Sugar()

		zap.ReplaceGlobals(zapLogger)

		zapLogger.Debug("initLogger: logger has been set globally")
	})
}

func Debug(args ...interface{}) {
	logger.Debug(args)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args)
}

func Info(args ...interface{}) {
	logger.Info(args)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args)
}

func Error(args ...interface{}) {
	logger.Info(args)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args)
}

func Panic(args ...interface{}) {
	logger.Panic(args)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args)
}

func Desugar() *zap.Logger {
	return logger.Desugar()
}

func Sync() error {
	return logger.Sync()
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
