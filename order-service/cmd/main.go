package main

import (
	"log"
	"net"
	"os"
	"time"

	"github.com/bekzxt/e-commerce/order-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/db"
	events "github.com/bekzxt/e-commerce/order-service/internal/infrastructure/events"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/repository"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/router"
	gr "github.com/bekzxt/e-commerce/order-service/internal/interfaces/grpc"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/http"
	orderpb "github.com/bekzxt/e-commerce/order-service/proto"
	reviewpb "github.com/bekzxt/e-commerce/order-service/proto_review"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
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
	var publisher *events.RabbitMQPublisher
	for {
		publisher, err = events.NewRabbitMQPublisher()
		if err == nil {
			log.Println("✅ Connected to RabbitMQ:", os.Getenv("RABBITMQ_URL"))
			break
		}

		log.Println("⚠️ RabbitMQ not ready yet, retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}
	defer publisher.Close()
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

	orderUseCase := usecase.NewOrderUseCase(orderRepo, orderItemRepo, publisher)
	reviewUseCase := usecase.NewReviewUseCase(reviewRepo)

	// HTTP handler
	orderHandler := http.NewOrderHandler(orderUseCase)

	// gRPC handlers
	orderGRPC := gr.NewOrderHandler(orderUseCase)
	reviewGRPC := gr.NewReviewHandler(reviewUseCase)

	// Start gRPC server
	go func() {
		lis, err := net.Listen("tcp", ":50052")
		if err != nil {
			log.Fatalf("❌ gRPC failed: %v", err)
		}
		grpcServer := grpc.NewServer()
		orderpb.RegisterOrderServiceServer(grpcServer, orderGRPC)
		reviewpb.RegisterReviewServiceServer(grpcServer, reviewGRPC)

		log.Println("✅ gRPC Order Service running on :50052")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("❌ Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server
	r := router.SetupRouter(*orderHandler)
	log.Println("✅ HTTP Order Service running on :8082")
	if err := r.Run(":8082"); err != nil {
		log.Fatalf("❌ Failed to start HTTP: %v", err)
	}
}
