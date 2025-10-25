package main

import (
	"github.com/bekzxt/e-commerce/order-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/db"
	events "github.com/bekzxt/e-commerce/order-service/internal/infrastructure/events"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/repository"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/router"
	gr "github.com/bekzxt/e-commerce/order-service/internal/interfaces/grpc"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/http"
	orderpb "github.com/bekzxt/e-commerce/order-service/proto"
	reviewpb "github.com/bekzxt/e-commerce/order-service/proto_review"
	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/joho/godotenv"
	ginprom "github.com/zsais/go-gin-prometheus"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	// Load env
	_ = godotenv.Load("../.env")

	// Connect DB
	database, err := db.ConnectPostgres()
	if err != nil {
		log.Fatalf("❌ Failed to connect to DB: %v", err)
	}
	defer database.Close()
	if err != nil {
		log.Fatalf("❌ Failed to connect RabbitMQ after retries: %v", err)
	}
	// Build dependencies
	orderRepo := repository.NewOrderRepository(database)
	orderItemRepo := repository.NewOrderItemRepository(database)
	reviewRepo := repository.NewReviewRepository(database)
	statusConsumer, cleanup, err := events.NewOrderStatusConsumer(orderRepo)
	if err != nil {
		log.Println("⚠️ OrderStatusConsumer init failed:", err)
	} else {
		if err := statusConsumer.Start(); err != nil {
			log.Println("⚠️ OrderStatusConsumer start failed:", err)
		} else {
			log.Println("✅ OrderStatusConsumer started")
		}
		defer cleanup()
	}
	invClient := infrastructure.NewInventoryClient("http://inventory-service:8081")
	orderUseCase := usecase.NewOrderUseCase(orderRepo, orderItemRepo, invClient)
	reviewUseCase := usecase.NewReviewUseCase(reviewRepo)

	// HTTP handler
	orderHandler := http.NewOrderHandler(orderUseCase)

	// gRPC handlers
	orderGRPC := gr.NewOrderHandler(orderUseCase)
	reviewGRPC := gr.NewReviewHandler(reviewUseCase)

	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
	)

	// Регистрируем сервисы
	orderpb.RegisterOrderServiceServer(s, orderGRPC)
	reviewpb.RegisterReviewServiceServer(s, reviewGRPC)

	// Регистрируем метрики Prometheus
	grpc_prometheus.Register(s)

	// Стартуем gRPC сервер
	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("❌ Failed to listen gRPC: %v", err)
		}

		log.Println("✅ gRPC Server running on :50052")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("❌ Failed to serve gRPC: %v", err)
		}
	}()

	infrastructure.StartMetricsServer()
	// Start HTTP server
	r := gin.Default()

	// ✅ Добавляем HTTP метрики
	p := ginprom.NewPrometheus("order_service")
	p.Use(r)

	// ✅ Регистрируем маршруты
	router.SetupRoutes(r, *orderHandler)

	log.Println("✅ HTTP Order Service running on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("❌ Failed to start HTTP: %v", err)
	}
}
