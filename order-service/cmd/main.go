package main

import (
	"github.com/bekzxt/e-commerce/order-service/internal/application/usecase"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/db"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/repository"
	"github.com/bekzxt/e-commerce/order-service/internal/infrastructure/router"
	gr "github.com/bekzxt/e-commerce/order-service/internal/interfaces/grpc"
	"github.com/bekzxt/e-commerce/order-service/internal/interfaces/http"
	orderpb "github.com/bekzxt/e-commerce/order-service/proto"
	reviewpb "github.com/bekzxt/e-commerce/order-service/proto_review"
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
	database, err := db.ConnectPostgres()
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer database.Close()

	r := gin.Default()

	orderRepo := repository.NewOrderRepository(database)
	orderItemRepo := repository.NewOrderItemRepository(database)
	reviewRepo := repository.NewReviewRepository(database)
	orderUseCase := usecase.NewOrderUseCase(orderRepo, orderItemRepo)
	orderHandler := http.NewOrderHandler(orderUseCase)
	reviewUseCase := usecase.NewReviewUseCase(reviewRepo)
	orderGRPChandler := gr.NewOrderHandler(orderUseCase)
	reviewGRPChandler := gr.NewReviewHandler(reviewUseCase)
	s := grpc.NewServer()
	orderpb.RegisterOrderServiceServer(s, orderGRPChandler)
	reviewpb.RegisterReviewServiceServer(s, reviewGRPChandler)

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Order Service running on :50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	r = router.SetupRouter(*orderHandler)
	r.Run(":8082")
}
