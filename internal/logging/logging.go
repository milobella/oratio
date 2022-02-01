package logging

import (
	"github.com/labstack/echo/v4"
	"github.com/milobella/oratio/pkg/logrusecho"
)

func ApplyMiddleware(server *echo.Echo) {
	builder := logrusecho.NewLogrusMiddlewareBuilder().
		ActivatedRequestData([]string{"request", "method"}).
		ActivatedResponseData([]string{"status"})
	server.Use(builder.Build().Handle)
}
