package logger

import (
	"os"

	"github.com/Karagar/cyprusExercise/pkg/utils"
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

// getLogger - singleton SugarLogger wrapper
func New() *zap.SugaredLogger {
	if logger != nil {
		return logger
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "DEBUG"
	}

	if logLevel == "DEBUG" {
		logger = getDevelopmentLog()
	} else {
		logger = getProductionLog()
	}
	return logger
}

// getProductionLog - wrapper for prod logger initiation
func getProductionLog() *zap.SugaredLogger {
	log, err := zap.NewProduction()
	utils.PanicOnErr(err)
	return log.Sugar()
}

// getDevelopmentLog - wrapper for dev logger initiation
func getDevelopmentLog() *zap.SugaredLogger {
	log, err := zap.NewDevelopment()
	utils.PanicOnErr(err)
	return log.Sugar()
}
