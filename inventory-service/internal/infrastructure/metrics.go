package infrastructure

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginprom "github.com/zsais/go-gin-prometheus"
	"google.golang.org/grpc"
)

const MetricsAddr = ":2113"

func AttachGRPCMetrics(s *grpc.Server) {
	grpc_prometheus.Register(s)
	grpc_prometheus.EnableHandlingTimeHistogram()
}

func NewGRPCServerWithMetrics(opts ...grpc.ServerOption) *grpc.Server {
	u := grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor)
	st := grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor)
	all := []grpc.ServerOption{u, st}
	all = append(all, opts...)
	s := grpc.NewServer(all...)
	AttachGRPCMetrics(s)
	return s
}

// Правильно добавляем HTTP метрики для Gin
func AttachGinMetrics(r *gin.Engine, serviceName string) {
	p := ginprom.NewPrometheus(serviceName)
	p.Use(r)
}

func StartMetricsServer() {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Printf("📈 Metrics server running on %s\n", MetricsAddr)
		if err := http.ListenAndServe(MetricsAddr, mux); err != nil {
			log.Printf("❌ Metrics server error: %v", err)
		}
	}()
}
