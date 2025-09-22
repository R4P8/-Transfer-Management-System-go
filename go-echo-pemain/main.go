package main

import (
	"context"
	"log"
	"net"

	clubpb "Pemain/Club/proto/club"
	pemainpb "Pemain/Pemain/proto/pemain"
	"Pemain/config"
	grpcserver "Pemain/grpcserver"
	"Pemain/repository"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(" .env file not found, using system environment")
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()
	db, err := config.DatabaseConnection(ctx)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()
	repo := repository.NewPemainRepository(db)

	//  connect ke service club
	clubConn, err := grpc.Dial("localhost:50053", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect club service: %v", err)
	}
	defer clubConn.Close()
	clubClient := clubpb.NewClubServiceClient(clubConn)

	grpcServer := grpc.NewServer()

	pemainServer := grpcserver.NewPemainServer(repo, clubClient)
	pemainpb.RegisterPemainServiceServer(grpcServer, pemainServer)

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
