package tracing

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer(serviceName, jaegerEndpoint string) (func(context.Context) error, error) {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return nil, fmt.Errorf("failed to create jaeger exporter: %w", err)
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	provider := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exporter),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
	)

	otel.SetTracerProvider(provider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func(ctx context.Context) error {
		return provider.Shutdown(ctx)
	}, nil
}

func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

func ExtractTraceInfoFromContext(ctx context.Context) (traceID, spanID string) {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return "", ""
	}

	spanContext := span.SpanContext()
	if spanContext.HasTraceID() {
		traceID = spanContext.TraceID().String()
	}
	if spanContext.HasSpanID() {
		spanID = spanContext.SpanID().String()
	}

	return traceID, spanID
}
