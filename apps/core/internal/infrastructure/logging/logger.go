package logging

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Handle struct {
	Logger *zap.Logger
	file   *lumberjack.Logger
}

func New(cfg config.LoggingConfig) (*Handle, error) {
	if err := os.MkdirAll(filepath.Dir(cfg.Path), 0o700); err != nil {
		return nil, fmt.Errorf("create log directory: %w", err)
	}

	var level zapcore.Level
	if err := level.Set(cfg.Level); err != nil {
		return nil, fmt.Errorf("parse log level: %w", err)
	}

	file := &lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSizeMB,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAgeDays,
		Compress:   cfg.Compress,
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder

	fileCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(file),
		level,
	)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.Lock(os.Stderr),
		level,
	)

	logger := zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller())
	return &Handle{Logger: logger, file: file}, nil
}

func (h *Handle) Close() error {
	if h == nil {
		return nil
	}

	var closeErrors []error
	if h.Logger != nil {
		if err := h.Logger.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			closeErrors = append(closeErrors, fmt.Errorf("sync logger: %w", err))
		}
	}
	if h.file != nil {
		if err := h.file.Close(); err != nil {
			closeErrors = append(closeErrors, fmt.Errorf("close log file: %w", err))
		}
	}

	return errors.Join(closeErrors...)
}
