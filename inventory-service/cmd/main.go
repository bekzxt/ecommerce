package main

import (
	_ "database/sql"
	"github.com/bekzxt/e-commerce/inventory-service/internal/infrastructure"
	"log"
	"net"
	"time"

	"github.com/bekzxt/e-commerce/inventory-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/inventory-service/internal/infrastructure/db"
	"github.com/bekzxt/e-commerce/inventory-service/internal/infrastructure/events"
	"github.com/bekzxt/e-commerce/inventory-service/internal/infrastructure/repository"
	grpchandler "github.com/bekzxt/e-commerce/inventory-service/internal/interfaces/grpc"
	handler "github.com/bekzxt/e-commerce/inventory-service/internal/interfaces/http"
	pb "github.com/bekzxt/e-commerce/inventory-service/proto"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load("../.env")

	// Connect DB
	database, err := db.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer database.Close()

	// Start Metrics Server
	infrastructure.StartMetricsServer()

	// Init gRPC with Prometheus
	s := infrastructure.NewGRPCServerWithMetrics()

	// Register repositories
	productRepo := repository.NewProductRepo(database)
	productUC := usecase.NewProductUseCase(productRepo)
	productHandler := grpchandler.NewProductHandler(productUC)
	pb.RegisterInventoryServiceServer(s, productHandler)

	// RabbitMQ Publisher
	var pub *events.RabbitMQPublisher
	for {
		p, err := events.NewRabbitMQPublisher()
		if err == nil {
			pub = p
			log.Println("✅ Connected to RabbitMQ Publisher (inventory)")
			break
		}
		log.Println("⚠️ RabbitMQ publisher retry in 3s...")
		time.Sleep(3 * time.Second)
	}
	defer pub.Close()

	// RabbitMQ Consumer
	var consumer *events.RabbitMQConsumer
	for {
		consumer, err = events.NewRabbitMQConsumer()
		if err == nil {
			log.Println("✅ Connected to RabbitMQ Consumer")
			break
		}
		log.Println("⚠️ RabbitMQ consumer retry in 3s...")
		time.Sleep(3 * time.Second)
	}
	defer consumer.Close()

	orderConsumer := events.NewOrderCreatedConsumer(consumer.Channel(), productUC, pub)
	go func() {
		if err := orderConsumer.Consume(); err != nil {
			log.Fatal("❌ Order consumer failed:", err)
		}
	}()

	// Start gRPC
	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("❌ failed to listen: %v", err)
		}
		log.Println("🚀 gRPC Inventory Service running on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC serve error: %v", err)
		}
	}()

	// Start HTTP REST API + Prometheus Gin
	r := gin.Default()
	infrastructure.AttachGinMetrics(r, "inventory_service")

	productHandlerHTTP := handler.NewProductHandler(productUC)
	discountRepo := repository.NewDiscountRepo(database)
	discountUC := usecase.NewDiscountUseCase(discountRepo)
	discountHandlerHTTP := handler.NewDiscountHandler(discountUC)

	productHandlerHTTP.RegisterRoutes(r)
	discountHandlerHTTP.RegisterRoutes(r)

	log.Println("🌍 REST Inventory API running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("REST API error: %v", err)
	}
}
