package logger

import (
	"go.uber.org/zap"
)

var Zap *zap.SugaredLogger

func InitializeLogger() {
	Zap = zap.NewExample().Sugar()
	Zap.Info("Logger has been initialized!")
}
