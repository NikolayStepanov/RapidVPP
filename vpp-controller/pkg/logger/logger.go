package logger

import (
	"context"
	"os"

	"github.com/NikolayStepanov/RapidVPP/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"

	"github.com/fsnotify/fsnotify"
)

var global *zap.Logger

func ObserveLoggerConfigFile(ctx context.Context, config *config.Config) *fsnotify.Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Fatal("observe logger create watcher", zap.Error(err))
	}

	err = watcher.Add(config.Logger.NameConfigFile)
	if err != nil {
		Fatal("observe logger level", zap.Error(err))
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					if global != nil {
						_ = global.Sync()
					}
					InitLogger(&config.Logger)
				}
			case err := <-watcher.Errors:
				if err != nil {
					Error("error while tracking changes in the logger configuration file", zap.Error(err))
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return watcher
}

func InitLogger(configLogger *config.LoggerConfig) {
	configFile, err := os.Open(configLogger.NameConfigFile)
	if err != nil {
		Fatal("init logger open file", zap.Error(err))
	}
	defer configFile.Close()

	decoder := yaml.NewDecoder(configFile)
	if err = decoder.Decode(&configLogger); err != nil {
		Error("init logger yaml decode", zap.Error(err))
	}

	level := zapcore.InvalidLevel
	if err = level.UnmarshalText([]byte(configLogger.Level)); err != nil {
		Fatal("init logger yaml decode level", zap.Error(err))
	}
	var timeEncoder zapcore.TimeEncoder
	if err = timeEncoder.UnmarshalText([]byte(configLogger.EncoderTime)); err != nil {
		Fatal("init logger yaml decode time encoder", zap.Error(err))
	}

	configZap := zap.NewProductionConfig()
	levelZap := zap.NewAtomicLevelAt(level)
	configZap.OutputPaths = configLogger.OutputPaths
	configZap.Level = levelZap
	configZap.Level.String()
	configZap.DisableCaller = true
	configZap.DisableStacktrace = true
	configZap.EncoderConfig.EncodeTime = timeEncoder
	newLogger, err := configZap.Build()
	if err != nil {
		Fatal("init logger build constructs", zap.Error(err))
	}
	global = newLogger
	global.WithOptions()
}

func Sync() error {
	return global.Sync()
}

func With(fields ...zap.Field) *zap.Logger {
	return global.With(fields...)
}

func Debug(msg string, fields ...zap.Field) {
	global.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	global.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	global.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	global.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	global.Fatal(msg, fields...)
}
