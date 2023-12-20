package pkg

import (
	"os"

	"go.uber.org/zap"
)

const (
	AccountIDKey = "account_id"
	RouteNameKey = "route_name"
	ReqBodyKey   = "request_body"
)

func NewLogger() *zap.SugaredLogger {
	var initLogger *zap.Logger
	initLogger, _ = zap.NewProduction()

	if os.Getenv("IS_DEV") == "1" {
		initLogger, _ = zap.NewDevelopment()
	}

	return initLogger.Sugar()
}
