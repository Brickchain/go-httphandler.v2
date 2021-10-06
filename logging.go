package httphandler

import (
	"context"

	"github.com/IpsoVeritas/logger"
)

func LoggerForContext(ctx context.Context) *logger.Entry {
	reqID, _ := ctx.Value(RequestIDKey).(string)
	return logger.WithField("request-id", reqID)
}
