package logger

import (
	"go.uber.org/zap"
)

func GetLogger(label string) *zap.Logger {
	config := zap.NewProductionConfig()
	config.InitialFields = make(map[string]interface{})
	config.InitialFields["label"] = label
	if logger, err := config.Build(); err != nil {
		panic(err)
	} else {
		return logger
	}
}
