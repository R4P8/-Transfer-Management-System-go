package client

import (
	"context"
	"fmt"

	playerpb "Pemain/Pemain/proto/pemain" // hasil generate dari player.proto

	"google.golang.org/grpc"
)

type PlayerClient struct {
	conn   *grpc.ClientConn
	client playerpb.PlayerServiceClient
}

// NewPlayerClient membuat koneksi ke service pemain
func NewPlayerClient(addr string) (*PlayerClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect player service: %w", err)
	}

	client := playerpb.NewPlayerServiceClient(conn)

	return &PlayerClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close untuk menutup koneksi gRPC
func (c *PlayerClient) Close() error {
	return c.conn.Close()
}

// GetPlayerByID memanggil service Player
func (c *PlayerClient) GetPlayerByID(ctx context.Context, id int64) (*playerpb.Player, error) {
	res, err := c.client.GetPlayerByID(ctx, &playerpb.GetPlayerRequest{Id: id})
	if err != nil {
		return nil, err
	}
	return res.Player, nil
}
