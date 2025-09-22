package client

import (
	"context"
	"fmt"

	clubpb "Pemain/Club/proto/club" // hasil generate dari club.proto

	"google.golang.org/grpc"
)

type ClubClient struct {
	conn   *grpc.ClientConn
	client clubpb.ClubServiceClient
}

// NewClubClient membuat koneksi ke service club
func NewClubClient(addr string) (*ClubClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to connect club service: %w", err)
	}

	client := clubpb.NewClubServiceClient(conn)

	return &ClubClient{
		conn:   conn,
		client: client,
	}, nil
}

// Close untuk menutup koneksi gRPC
func (c *ClubClient) Close() error {
	return c.conn.Close()
}

// GetClubByID memanggil service Club
func (c *ClubClient) GetClubByID(ctx context.Context, id int64) (*clubpb.Club, error) {
	res, err := c.client.GetClubByID(ctx, &clubpb.GetClubByIDRequest{IdClub: id})
	if err != nil {
		return nil, err
	}
	return res.Club, nil
}
