package logging

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gofurry/corvus-studio/apps/core/internal/infrastructure/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Handle struct {
	Logger *zap.Logger
	file   *lumberjack.Logger
}

type writeOnly struct {
	io.Writer
}

func New(cfg config.LoggingConfig) (*Handle, error) {
	return newWithConsole(cfg, os.Stderr)
}

func newWithConsole(cfg config.LoggingConfig, console io.Writer) (*Handle, error) {
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
		zapcore.Lock(zapcore.AddSync(writeOnly{Writer: console})),
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
		if err := h.Logger.Sync(); err != nil {
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
