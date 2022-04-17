package utils

import (
	"go.uber.org/zap"
)

// getProductionLog - wrapper for prod logger initiation
func getProductionLog() *zap.SugaredLogger {
	log, err := zap.NewProduction()
	panicOnErr(err)
	return log.Sugar()
}

// getDevelopmentLog - wrapper for dev logger initiation
func getDevelopmentLog() *zap.SugaredLogger {
	log, err := zap.NewDevelopment()
	panicOnErr(err)
	return log.Sugar()
}
