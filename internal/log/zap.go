package log

import (
	"context"
	"go.uber.org/zap"
	l "log"
	"os"
)

var log *zap.SugaredLogger

const (
	logEnvPreset     = "LOG_PRESET"
	productionPreset = "prod"
	devPreset        = "dev"
)

type logConfig struct {
	preset string
}

func InitLogger(_ context.Context) error {
	zapLog, err := selectLogger()
	if err != nil {
		return err
	}

	log = zapLog.Sugar()

	return nil
}

func selectLogger() (*zap.Logger, error) {
	switch newLogConfig().preset {
	case productionPreset:
		return zap.NewProduction()
	case devPreset:
		return zap.NewDevelopment()
	default:
		l.Println("unknown logger preset, using development preset")
		return zap.NewDevelopment()
	}
}

func newLogConfig() *logConfig {
	return &logConfig{preset: os.Getenv(logEnvPreset)}
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Debugf(template string, args ...interface{}) {
	log.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	log.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	log.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	log.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	log.Fatalf(template, args...)
}
