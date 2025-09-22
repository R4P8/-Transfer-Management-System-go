package main

import (
	"context"
	"log"
	"net"

	clubpb "Club/Club/proto/club"
	"Club/config"
	grpcserver "Club/grpcserver"
	"Club/repository"
	"Club/service"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(" .env file not found, using system environment")
	}

	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()
	db, err := config.DatabaseConnection(ctx)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()

	repo := repository.NewClubRepository(db)
	svc := service.NewClubService(repo)

	playerAddr := "localhost:50051"
	conn, err := grpc.Dial(playerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect pemain service: %v", err)
	}
	defer conn.Close()

	// bikin gRPC server transfer
	grpcServer := grpc.NewServer()
	clubServer := grpcserver.NewClubServer(svc)

	//  register club service
	clubpb.RegisterClubServiceServer(grpcServer, clubServer)

	log.Println("gRPC Transfer server listening on :50053")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
