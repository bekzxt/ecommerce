package main

import (
	_ "database/sql"
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
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("Error loading .env file")
	}
	database, err1 := db.ConnectPostgres()
	if err1 != nil {
		log.Fatalf("Failed to connect to DB: %v", err1)
	}
	defer database.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	productRepo := repository.NewProductRepo(database)
	productUC := usecase.NewProductUseCase(productRepo)

	var pub *events.RabbitMQPublisher
	for {
		p, err := events.NewRabbitMQPublisher()
		if err == nil {
			pub = p
			log.Println("✅ Connected to RabbitMQ Publisher (inventory):", os.Getenv("RABBITMQ_URL"))
			break
		}
		log.Println("⚠️ RabbitMQ not ready for publisher (inventory), retrying in 3s...", err)
		time.Sleep(3 * time.Second)
	}
	defer pub.Close()

	var consumer *events.RabbitMQConsumer
	for {
		consumer, err = events.NewRabbitMQConsumer()
		if err == nil {
			log.Println("✅ Connected to RabbitMQ Consumer:", os.Getenv("RABBITMQ_URL"))
			break
		}

		log.Println("⚠️ RabbitMQ not ready yet for Consumer, retrying in 3 seconds...")
		time.Sleep(3 * time.Second)
	}
	defer consumer.Close()
	orderConsumer := events.NewOrderCreatedConsumer(consumer.Channel(), productUC, pub)

	go func() {

		if err := orderConsumer.Consume(); err != nil {
			log.Fatal("❌ Failed to start order.created consumer:", err)
		}
	}()

	log.Println("✅ RabbitMQ order.created consumer started")

	productHandler := grpchandler.NewProductHandler(productUC)
	s := grpc.NewServer()
	pb.RegisterInventoryServiceServer(s, productHandler)
	discountRepo := repository.NewDiscountRepo(database)
	discountUC := usecase.NewDiscountUseCase(discountRepo)
	discountHandlerr := handler.NewDiscountHandler(discountUC)
	productUCgin := usecase.NewProductUseCase(productRepo)
	productHandlerr := handler.NewProductHandler(productUCgin)
	r := gin.Default()
	productHandlerr.RegisterRoutes(r)
	discountHandlerr.RegisterRoutes(r)
	log.Println("Inventory Service running on :50051")

	go func() {
		log.Println("🚀 Starting gRPC Inventory Service on :50051")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	log.Println("🌍 Starting REST Inventory API on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to run REST API: %v", err)
	}
}
