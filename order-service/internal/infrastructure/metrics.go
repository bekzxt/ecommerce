package infrastructure

import (
	"log"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

const MetricsAddr = ":2112"

// AttachGRPCMetrics регистрирует стандартные gRPC пром-метрики и возвращает опции для сервера.
func AttachGRPCMetrics(s *grpc.Server) {
	grpc_prometheus.Register(s) // регистрирует дефолтные серверные метрики на s
	grpc_prometheus.EnableHandlingTimeHistogram()
}

// NewGRPCServerWithMetrics создаёт gRPC-сервер c interceptors для метрик.
func NewGRPCServerWithMetrics(opts ...grpc.ServerOption) *grpc.Server {
	u := grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor)
	st := grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor)

	all := []grpc.ServerOption{u, st}
	all = append(all, opts...)
	s := grpc.NewServer(all...)
	AttachGRPCMetrics(s)
	return s
}

// StartMetricsServer поднимает /metrics на 2112
func StartMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Printf("📈 order-service metrics on %s\n", MetricsAddr)
		if err := http.ListenAndServe(MetricsAddr, mux); err != nil {
			log.Printf("metrics server error: %v", err)
		}
	}()
}
