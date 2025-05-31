package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"route256/cart/internal/api/middlewares"
	"route256/cart/internal/logging"
	"route256/cart/internal/metrics"
	"route256/cart/internal/models"
	"route256/cart/internal/tracing"

	"go.uber.org/zap"
)

type Router struct {
	address     string
	cartService CartService
	tracer      string
	jaegerURL   string
}

//go:generate minimock -i route256/cart/internal/api.CartService -o mocked/cart_service_mock.go -n CartServiceMock -p mocked

type CartService interface {
	AddProduct(ctx context.Context, userID, sku int64, count uint32) error
	DeleteProduct(ctx context.Context, userID, sku int64) error
	Cart(ctx context.Context, userID int64) (*models.Cart, error)
	Checkout(ctx context.Context, userID int64) (int64, error)
	ClearCart(userID int64) bool
}

func NewRouter(
	host, port string,
	cart CartService,
) *Router {
	return &Router{
		address:     fmt.Sprintf("%s:%s", host, port),
		cartService: cart,
		tracer:      "cart",
		jaegerURL:   "http://jaeger:14268/api/traces",
	}
}

func (r *Router) Start() error {
	cleanup, err := tracing.InitTracer(r.tracer, r.jaegerURL)
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}
	defer func() {
		if err = cleanup(context.Background()); err != nil {
			logging.Error("Failed to clean up tracer", zap.Error(err))
		}
	}()

	logging.InitLogger()
	defer func() {
		err = logging.Sync()
		if err != nil {
			logging.Error("Failed to sync logger", zap.Error(err))
		}
	}()

	mx := http.NewServeMux()

	mx.HandleFunc("POST /user/{user_id}/cart/{sku}", r.addProduct)
	mx.HandleFunc("DELETE /user/{user_id}/cart/{sku}", r.deleteProduct)
	mx.HandleFunc("DELETE /user/{user_id}/cart", r.clearCart)
	mx.HandleFunc("GET /user/{user_id}/cart", r.cart)
	mx.HandleFunc("POST /user/{user_id}/cart/checkout", r.checkout)

	mx.HandleFunc("GET /user/{user_id}/order/pay", r.cart)
	mx.HandleFunc("GET /user/{user_id}/order/cancel", r.cart)

	mx.Handle("GET /metrics", metrics.NewMetricsHandler())

	middleware := middlewares.NewLoggingMiddleware(mx)

	logging.Info("Starting router", zap.String("address", r.address))

	server := http.Server{
		Addr:              r.address,
		Handler:           middleware,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err = server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("can't start router: %w", err)
	}

	return nil
}

func ErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")

	response := struct {
		Message string `json:"message"`
	}{
		Message: message,
	}

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		fmt.Printf("can't encode response: %v; error message: %s", err, message)
	}
}
