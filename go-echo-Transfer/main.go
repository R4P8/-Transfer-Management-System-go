package main

import (
	"context"
	"log"
	"net"

	pemainpb "Transfer/Pemain/proto/pemain"
	"Transfer/Transfer/proto/transferpb"
	"Transfer/config"
	grpcserver "Transfer/grpcserver"
	"Transfer/repository"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println(" .env file not found, using system environment")
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	ctx := context.Background()
	db, err := config.DatabaseConnection(ctx)
	if err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}
	defer db.Close()
	repo := repository.NewTransferRepository(db)

	playerAddr := "localhost:50051"
	conn, err := grpc.Dial(playerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect pemain service: %v", err)
	}
	defer conn.Close()

	pemainClient := pemainpb.NewPemainServiceClient(conn)

	// bikin gRPC server transfer
	grpcServer := grpc.NewServer()
	transferServer := grpcserver.NewTransferServer(repo, pemainClient)

	//  register transfer service
	transferpb.RegisterTransferServiceServer(grpcServer, transferServer)

	log.Println("gRPC Transfer server listening on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
