package api

import (
	"context"
	"fmt"
	"time"

	"route256/loms/internal/logging"
	"route256/loms/internal/metrics"
	"route256/loms/internal/tracing"

	"go.opentelemetry.io/otel/attribute"
	otlpcodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	startTime := time.Now()

	method := info.FullMethod

	resp, err := handler(ctx, req)

	duration := time.Since(startTime)

	statusCode := status.Code(err)
	statusCodeInt := int(statusCode)

	metrics.ObserveGRPCRequest(method, statusCodeInt, duration)

	return resp, err
}

func TracingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	method := info.FullMethod

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		propagator := propagation.TraceContext{}
		ctx = propagator.Extract(ctx, propagation.HeaderCarrier(md))
	}

	ctx, span := tracing.Tracer("loms-grpc").Start(ctx, method)
	defer span.End()

	traceID, spanID := tracing.ExtractTraceInfoFromContext(ctx)

	ctx = context.WithValue(ctx, "trace_id", traceID)
	ctx = context.WithValue(ctx, "span_id", spanID)

	span.SetAttributes(
		attribute.String("grpc.method", method),
		attribute.String("grpc.request", fmt.Sprintf("%+v", req)),
	)

	logging.Info("Incoming gRPC request",
		zap.String("method", method),
		zap.Any("request", req),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)

	resp, err := handler(ctx, req)

	if err != nil {
		st, _ := status.FromError(err)
		span.SetAttributes(attribute.String("grpc.status_code", st.Code().String()))
		span.SetStatus(otlpcodes.Error, st.Message())
		span.RecordError(err)

		logging.Error("gRPC request failed",
			zap.String("method", method),
			zap.Error(err),
			zap.String("status_code", st.Code().String()),
			zap.String("trace_id", traceID),
		)
	} else {
		span.SetAttributes(attribute.String("grpc.status_code", codes.OK.String()))
		span.SetStatus(otlpcodes.Ok, "")

		logging.Info("gRPC request succeeded",
			zap.String("method", method),
			zap.Any("response", resp),
			zap.String("trace_id", traceID),
		)
	}

	return resp, err
}

func Validate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if v, ok := req.(interface{ ValidateAll() error }); ok {
		if err := v.ValidateAll(); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}
	return handler(ctx, req)
}
