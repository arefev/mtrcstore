package logger

import (
	"fmt"

	"go.uber.org/zap"
)

var Log *zap.Logger

func Init(level string) error {
    lvl, err := zap.ParseAtomicLevel(level)
    if err != nil {
        return fmt.Errorf("zap logger parse level failed: %w", err)
    }

    cfg := zap.NewProductionConfig()
    cfg.Level = lvl

    // создаём логер на основе конфигурации
    zl, err := cfg.Build()
    if err != nil {
        return fmt.Errorf("zap logger build from config failed: %w", err)
    }

    Log = zl
    return nil
}