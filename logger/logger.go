package logger

import (
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger     *zap.Logger
	globalLoggerOnce sync.Once
)

func InitLogger(c *Config) {
	if c == nil {
		c = DefaultConfig()
	}
	globalLoggerOnce.Do(func() {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			FunctionKey:    zapcore.OmitKey,
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		encoder := zapcore.NewJSONEncoder(encoderConfig)

		logLevel := getZapLevel(c.Level)
		levelEnabler := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= logLevel
		})

		var core zapcore.Core
		switch strings.ToLower(c.Output) {
		case "stdout":
			core = zapcore.NewCore(
				encoder,
				zapcore.AddSync(os.Stdout),
				levelEnabler,
			)
		case "stderr":
			core = zapcore.NewCore(
				encoder,
				zapcore.AddSync(os.Stderr),
				levelEnabler,
			)
		case "file":
			fileWriter := zapcore.AddSync(&lumberjack.Logger{
				Filename:   c.LogFile,
				MaxSize:    c.MaxSize,
				MaxBackups: c.MaxFiles,
				MaxAge:     c.MaxAge,
				Compress:   c.Compress,
			})
			core = zapcore.NewCore(
				encoder,
				fileWriter,
				levelEnabler,
			)
		default:
			core = zapcore.NewCore(
				encoder,
				zapcore.AddSync(os.Stdout),
				levelEnabler,
			)
		}

		globalLogger = zap.New(core,
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)
	})
}

func getZapLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "err", "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func Logger() *zap.Logger {
	if globalLogger == nil {
		InitLogger(nil)
	}
	return globalLogger
}
