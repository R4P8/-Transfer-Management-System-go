package service

import (
	"Transfer/entities"
	"Transfer/repository"
	"context"
	"errors"

	pemainpb "Transfer/Pemain/proto/pemain" // ini import gRPC client service pemain

	"google.golang.org/grpc"
)

// TransferService interface
type TransferService interface {
	CreateTransfer(ctx context.Context, t *entities.Transfer) error
	GetTransferByID(ctx context.Context, id int64) (*entities.Transfer, error)
	GetAllTransfers(ctx context.Context) ([]*entities.Transfer, error)
	UpdateTransfer(ctx context.Context, t *entities.Transfer) error
	DeleteTransfer(ctx context.Context, id int64) error
}

// transferService struct
type transferService struct {
	repo       repository.TransferRepository
	pemainConn pemainpb.PemainServiceClient
}

// Constructor
func NewTransferService(repo repository.TransferRepository, playerServiceAddr string) (TransferService, error) {
	// connect ke service player via gRPC
	conn, err := grpc.Dial(playerServiceAddr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pemainpb.NewPemainServiceClient(conn)
	return &transferService{
		repo:       repo,
		pemainConn: client,
	}, nil
}

// Create transfer hanya jika status pemain sesuai rule
func (s *transferService) CreateTransfer(ctx context.Context, t *entities.Transfer) error {
	// cek status pemain dari service player
	res, err := s.pemainConn.GetPemainByID(ctx, &pemainpb.GetPemainByIDRequest{IdPemain: t.PlayerID})
	if err != nil {
		return err
	}

	allowed := map[string]bool{
		"free_agent":    true,
		"loan_list":     true,
		"transfer_list": true,
	}

	if !allowed[res.Pemain.StatusPemain] {
		return errors.New("transfer hanya bisa dilakukan jika status pemain free_agent, loan_list, atau transfer_list")
	}

	return s.repo.Create(ctx, t)
}

func (s *transferService) GetTransferByID(ctx context.Context, id int64) (*entities.Transfer, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *transferService) GetAllTransfers(ctx context.Context) ([]*entities.Transfer, error) {
	return s.repo.GetAll(ctx)
}

func (s *transferService) UpdateTransfer(ctx context.Context, t *entities.Transfer) error {
	return s.repo.Update(ctx, t)
}

func (s *transferService) DeleteTransfer(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
