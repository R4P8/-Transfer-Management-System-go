package grpcserver

import (
	"context"
	"errors"

	pemainpb "Transfer/Pemain/proto/pemain"
	"Transfer/Transfer/proto/transferpb"
	"Transfer/entities"
	"Transfer/repository"
	// import client dari service pemain
)

type TransferServer struct {
	transferpb.UnimplementedTransferServiceServer
	repo         repository.TransferRepository
	pemainClient pemainpb.PemainServiceClient
}

// Constructor
func NewTransferServer(repo repository.TransferRepository, pemainClient pemainpb.PemainServiceClient) *TransferServer {
	return &TransferServer{
		repo:         repo,
		pemainClient: pemainClient,
	}
}

// CreateTransfer -> cek status pemain dulu baru buat transfer
func (s *TransferServer) CreateTransfer(ctx context.Context, req *transferpb.CreateTransferRequest) (*transferpb.TransferResponse, error) {
	// Panggil service pemain untuk ambil data pemain
	playerResp, err := s.pemainClient.GetPemainByID(ctx, &pemainpb.GetPemainByIDRequest{
		IdPemain: req.PlayerId,
	})
	if err != nil {
		return nil, errors.New("failed to fetch player from pemain service")
	}
	player := playerResp.Pemain

	// Validasi status
	if player.StatusPemain != "free_agent" && player.StatusPemain != "loan_list" && player.StatusPemain != "transfer_list" {
		return nil, errors.New("player not eligible for transfer")
	}

	// Buat entity transfer
	t := &entities.Transfer{
		PlayerID:     player.IdPemain,
		OldClubID:    player.ClubId,
		NewClubID:    req.NewClubId,
		TransferFee:  req.TransferFee,
		TransferDate: req.TransferDate,
		Status:       "pending",
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}

	return &transferpb.TransferResponse{
		Transfer: &transferpb.Transfer{
			IdTransfer:   t.IDTransfer,
			PlayerId:     t.PlayerID,
			OldClubId:    t.OldClubID,
			NewClubId:    t.NewClubID,
			TransferFee:  t.TransferFee,
			TransferDate: t.TransferDate,
			Status:       t.Status,
		},
	}, nil
}

// GetTransfer
func (s *TransferServer) GetTransfer(ctx context.Context, req *transferpb.GetTransferRequest) (*transferpb.TransferResponse, error) {
	t, err := s.repo.GetByID(ctx, req.IdTransfer)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, errors.New("transfer not found")
	}

	return &transferpb.TransferResponse{
		Transfer: &transferpb.Transfer{
			IdTransfer:   t.IDTransfer,
			PlayerId:     t.PlayerID,
			OldClubId:    t.OldClubID,
			NewClubId:    t.NewClubID,
			TransferFee:  t.TransferFee,
			TransferDate: t.TransferDate,
			Status:       t.Status,
		},
	}, nil
}

// ListTransfers
func (s *TransferServer) ListTransfers(ctx context.Context, req *transferpb.ListTransfersRequest) (*transferpb.ListTransfersResponse, error) {
	transfers, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var pbTransfers []*transferpb.Transfer
	for _, t := range transfers {
		pbTransfers = append(pbTransfers, &transferpb.Transfer{
			IdTransfer:   t.IDTransfer,
			PlayerId:     t.PlayerID,
			OldClubId:    t.OldClubID,
			NewClubId:    t.NewClubID,
			TransferFee:  t.TransferFee,
			TransferDate: t.TransferDate,
			Status:       t.Status,
		})
	}

	return &transferpb.ListTransfersResponse{Transfers: pbTransfers}, nil
}

// UpdateTransfer
func (s *TransferServer) UpdateTransfer(ctx context.Context, req *transferpb.UpdateTransferRequest) (*transferpb.TransferResponse, error) {

	t := &entities.Transfer{
		IDTransfer:  req.IdTransfer,
		NewClubID:   req.NewClubId,
		TransferFee: req.TransferFee,
		Status:      req.Status,
	}

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return s.GetTransfer(ctx, &transferpb.GetTransferRequest{IdTransfer: req.IdTransfer})
}

// DeleteTransfer
func (s *TransferServer) DeleteTransfer(ctx context.Context, req *transferpb.DeleteTransferRequest) (*transferpb.DeleteTransferResponse, error) {
	if err := s.repo.Delete(ctx, req.IdTransfer); err != nil {
		return nil, err
	}
	return &transferpb.DeleteTransferResponse{Success: true}, nil
}
