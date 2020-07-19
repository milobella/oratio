package logging

import (
	"github.com/milobella/oratio/pkg/logrusecho"
)

func InitializeLoggingMiddleware() logrusecho.LogrusMiddleware{
	return logrusecho.NewLogrusMiddlewareBuilder().ActivatedRequestData(
		[]string{"request", "method"}).ActivatedResponseData([]string{"status"}).Build()
}
