package tracing

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/milobella/oratio/internal/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

func InitGlobalTracer(conf config.TracingConfig) (func(), error) {
	// Create the Jaeger exporter
	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(conf.JaegerAgentHostName),
		jaeger.WithAgentPort(fmt.Sprintf("%d", conf.JaegerAgentPort)),
	))
	if err != nil {
		return nil, err
	}

	resources := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("milobella"),
		semconv.ServiceVersionKey.String("1.0.0"),
		semconv.ServiceInstanceIDKey.String("oratio"),
	)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resources),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(b3.New())

	return func() {
		if shutdownErr := tp.Shutdown(context.Background()); shutdownErr != nil {
			logrus.Printf("Error shutting down tracer provider: %v", shutdownErr)
		}
	}, nil
}

func ApplyMiddleware(server *echo.Echo, conf config.TracingConfig) {
	server.Use(otelecho.Middleware(conf.ServiceName))
}
