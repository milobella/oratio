package logging

import "gitlab.milobella.com/milobella/oratio/pkg/logrushttp"

func InitializeLoggingMiddleware() logrushttp.LogrusMiddleware{
	return logrushttp.NewLogrusMiddlewareBuilder().ActivatedRequestData(
		[]string{"request", "method"}).ActivatedResponseData([]string{"status"}).Build()
}