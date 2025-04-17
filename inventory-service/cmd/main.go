package main

import (
	_ "database/sql"
	"github.com/bekzxt/e-commerce/inventory-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/inventory-service/internal/infrastructure/db"
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
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
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

	log.Println(" running on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
