package logger

import (
	"go.uber.org/zap"
)

type field = zap.Field

func Err(err error) field { return zap.Error(err) }

func New(env string) *zap.Logger {
	if env == "production" {
		l, _ := zap.NewProduction()
		return l
	}
	l, _ := zap.NewDevelopment()
	return l
} 