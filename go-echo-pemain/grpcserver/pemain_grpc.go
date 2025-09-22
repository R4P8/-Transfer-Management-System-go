package servergrpc

import (
	"context"
	"errors"

	clubpb "Pemain/Club/proto/club"
	pemainpb "Pemain/Pemain/proto/pemain"
	"Pemain/entities"
	"Pemain/repository"
)

// PemainServer implementasi gRPC
type PemainServer struct {
	pemainpb.UnimplementedPemainServiceServer
	repo       repository.PemainRepository
	clubClient clubpb.ClubServiceClient
}

// Constructor
func NewPemainServer(repo repository.PemainRepository, clubClient clubpb.ClubServiceClient) *PemainServer {
	return &PemainServer{
		repo:       repo,
		clubClient: clubClient,
	}
}

// ======================
// Implementasi RPC
// ======================

// Create pemain baru
func (s *PemainServer) CreatePemain(ctx context.Context, req *pemainpb.CreatePemainRequest) (*pemainpb.PemainResponse, error) {
	newP := &entities.Pemain{
		Name:        req.Name,
		ClubID:      &req.ClubId,
		MarketValue: req.MarketValue,
		Status:      req.StatusPemain,
	}

	if err := s.repo.Create(ctx, newP); err != nil {
		return nil, err
	}

	return &pemainpb.PemainResponse{
		Pemain: &pemainpb.Pemain{
			IdPemain:     newP.IDPemain,
			Name:         newP.Name,
			ClubId:       *newP.ClubID,
			MarketValue:  newP.MarketValue,
			StatusPemain: newP.Status,
		},
	}, nil
}

// Get pemain by ID
func (s *PemainServer) GetPemainByID(ctx context.Context, req *pemainpb.GetPemainByIDRequest) (*pemainpb.PemainResponse, error) {
	p, err := s.repo.GetByID(ctx, req.IdPemain)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, errors.New("pemain not found")
	}

	clubID := int64(0)
	if p.ClubID != nil {
		clubID = *p.ClubID
	}

	return &pemainpb.PemainResponse{
		Pemain: &pemainpb.Pemain{
			IdPemain:     p.IDPemain,
			Name:         p.Name,
			ClubId:       clubID,
			MarketValue:  p.MarketValue,
			StatusPemain: p.Status,
		},
	}, nil
}

// Get semua pemain
func (s *PemainServer) GetAllPemain(ctx context.Context, req *pemainpb.Empty) (*pemainpb.ListPemainResponse, error) {
	players, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var list []*pemainpb.Pemain
	for _, p := range players {
		clubID := int64(0)
		if p.ClubID != nil {
			clubID = *p.ClubID
		}
		list = append(list, &pemainpb.Pemain{
			IdPemain:     p.IDPemain,
			Name:         p.Name,
			ClubId:       clubID,
			MarketValue:  p.MarketValue,
			StatusPemain: p.Status,
		})
	}

	return &pemainpb.ListPemainResponse{Pemain: list}, nil
}

// Update oleh pemain sendiri (hanya status)
func (s *PemainServer) UpdateByPlayer(ctx context.Context, req *pemainpb.UpdateByPlayerRequest) (*pemainpb.PemainResponse, error) {
	if err := s.repo.UpdateByPlayer(ctx, req.IdPemain, req.StatusPemain); err != nil {
		return nil, err
	}
	return s.GetPemainByID(ctx, &pemainpb.GetPemainByIDRequest{IdPemain: req.IdPemain})

}

// Update oleh club (boleh club_id, market_value, status)
func (s *PemainServer) UpdateByClub(ctx context.Context, req *pemainpb.UpdateByClubRequest) (*pemainpb.PemainResponse, error) {
	err := s.repo.UpdateByClub(ctx, req.IdPemain, req.ClubId, req.MarketValue, req.StatusPemain)
	if err != nil {
		return nil, err
	}

	// Ambil data terbaru
	p, err := s.repo.GetByID(ctx, req.IdPemain)
	if err != nil {
		return nil, err
	}

	clubID := int64(0)
	if p.ClubID != nil {
		clubID = *p.ClubID
	}

	return &pemainpb.PemainResponse{
		Pemain: &pemainpb.Pemain{
			IdPemain:     p.IDPemain,
			Name:         p.Name,
			ClubId:       clubID,
			MarketValue:  p.MarketValue,
			StatusPemain: p.Status,
		},
	}, nil

	return s.GetPemainByID(ctx, &pemainpb.GetPemainByIDRequest{IdPemain: req.IdPemain})
}

// Delete pemain
func (s *PemainServer) DeletePemain(ctx context.Context, req *pemainpb.DeletePemainRequest) (*pemainpb.Empty, error) {
	if err := s.repo.Delete(ctx, req.IdPemain); err != nil {
		return nil, err
	}
	return &pemainpb.Empty{}, nil
}

// Get pemain berdasarkan club_id (dipanggil oleh service Club)
func (s *PemainServer) GetPemainByClub(ctx context.Context, req *pemainpb.GetPemainByClubRequest) (*pemainpb.ListPemainResponse, error) {
	players, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var list []*pemainpb.Pemain
	for _, p := range players {
		if p.ClubID != nil && *p.ClubID == req.ClubId {
			list = append(list, &pemainpb.Pemain{
				IdPemain:     p.IDPemain,
				Name:         p.Name,
				ClubId:       *p.ClubID,
				MarketValue:  p.MarketValue,
				StatusPemain: p.Status,
			})
		}
	}

	return &pemainpb.ListPemainResponse{Pemain: list}, nil
}
