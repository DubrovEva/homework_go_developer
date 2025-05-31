package middlewares

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"route256/cart/internal/logging"
	"route256/cart/internal/metrics"
	"route256/cart/internal/tracing"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ctxKey string

const (
	traceIDCtxKey ctxKey = "trace_id"
	spanIDCtxKey  ctxKey = "span_id"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

type LoggingMiddleware struct {
	handler http.Handler
	tracer  trace.Tracer
}

func NewLoggingMiddleware(handler http.Handler) http.Handler {
	return &LoggingMiddleware{
		handler: handler,
		tracer:  tracing.Tracer("cart-http"),
	}
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	ctx := r.Context()
	propagator := propagation.TraceContext{}
	ctx = propagator.Extract(ctx, propagation.HeaderCarrier(r.Header))

	ctx, span := m.tracer.Start(ctx, "HTTP "+r.Method+" "+r.URL.Path)
	defer span.End()

	traceID, spanID := tracing.ExtractTraceInfoFromContext(ctx)
	ctx = context.WithValue(ctx, traceIDCtxKey, traceID)
	ctx = context.WithValue(ctx, spanIDCtxKey, spanID)

	span.SetAttributes(
		attribute.String("http.method", r.Method),
		attribute.String("http.url", r.URL.String()),
		attribute.String("http.host", r.Host),
	)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logging.Error("Failed to read request body",
			zap.Error(err),
			zap.String("method", r.Method),
			zap.String("url", r.URL.String()),
		)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	rw := newResponseWriter(w)

	logging.Info("Incoming request",
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.String("body", string(body)),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)

	m.handler.ServeHTTP(rw, r.WithContext(ctx))

	duration := time.Since(startTime)

	metrics.ObserveHTTPRequest(r.Method, r.URL.Path, rw.statusCode, duration)

	span.SetAttributes(
		attribute.Int("http.status_code", rw.statusCode),
		attribute.Int64("http.response_time_ms", duration.Milliseconds()),
	)

	if rw.statusCode >= 400 {
		span.SetStatus(codes.Error, http.StatusText(rw.statusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}

	logging.Info("Outgoing response",
		zap.String("method", r.Method),
		zap.String("url", r.URL.String()),
		zap.Int("status", rw.statusCode),
		zap.Duration("duration", duration),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)
}
