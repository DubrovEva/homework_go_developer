package api

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"route256/loms/internal/logging"
	"route256/loms/internal/metrics"
	"route256/loms/internal/tracing"

	loms "route256/loms/pkg/api/loms/v1"

	"go.uber.org/zap"
)

type Gateway struct {
	lomsService *Loms

	grpcPort  int64
	httpPort  int64
	tracer    string
	jaegerURL string
}

func NewGateway(lomsService *Loms, grpcPort, httpPort int64) *Gateway {
	return &Gateway{
		lomsService: lomsService,
		grpcPort:    grpcPort,
		httpPort:    httpPort,
		tracer:      "loms",
		jaegerURL:   "http://jaeger:14268/api/traces",
	}
}

func (s *Gateway) Run(ctx context.Context) error {
	cleanup, err := tracing.InitTracer(s.tracer, s.jaegerURL)
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			logging.Error("Failed to clean up tracer", zap.Error(err))
		}
	}()

	logging.InitLogger()
	defer logging.Sync()

	list, err := net.Listen("tcp", fmt.Sprintf(":%d", s.grpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			MetricsInterceptor,
			TracingInterceptor,
			Validate,
		),
	)

	reflection.Register(grpcServer)
	loms.RegisterLomsServer(grpcServer, s.lomsService)

	logging.Info("Starting gRPC server", zap.Int64("port", s.grpcPort))
	go func() {
		if err = grpcServer.Serve(list); err != nil {
			logging.Error("Failed to serve gRPC server", zap.Error(err))
			panic(err)
		}
	}()

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", s.grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to create gRPC client: %v", err)
	}

	gatewayMux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(
			func(header string) (string, bool) {
				switch strings.ToLower(header) {
				case "x-auth":
					return header, true
				default:
					return header, false
				}
			},
		),
	)

	if err = loms.RegisterLomsHandler(ctx, gatewayMux, conn); err != nil {
		return fmt.Errorf("failed to register handler: %v", err)
	}

	httpMux := http.NewServeMux()

	httpMux.Handle("/", gatewayMux)

	httpMux.Handle("/metrics", metrics.NewMetricsHandler())

	gwServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", s.httpPort),
		Handler: httpMux,
	}

	logging.Info("Starting HTTP server", zap.Int64("port", s.httpPort))
	if err = gwServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to serve HTTP server: %v", err)
	}

	return nil
}
