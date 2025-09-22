package grpcserver

import (
	"context"

	clubpb "Club/Club/proto/club"
	"Club/entities"
	"Club/service"
)

type ClubServer struct {
	clubpb.UnimplementedClubServiceServer
	svc service.ClubService
}

func NewClubServer(svc service.ClubService) *ClubServer {
	return &ClubServer{svc: svc}
}

func (s *ClubServer) CreateClub(ctx context.Context, req *clubpb.CreateClubRequest) (*clubpb.ClubResponse, error) {
	club := &entities.Club{
		NameClub: req.NameClub,
		CityClub: req.CityClub,
		Budget:   req.Budget,
	}
	if err := s.svc.CreateClub(ctx, club); err != nil {
		return nil, err
	}
	return &clubpb.ClubResponse{
		Club: &clubpb.Club{
			IdClub:   club.IDClub,
			NameClub: club.NameClub,
			CityClub: club.CityClub,
			Budget:   club.Budget,
		},
	}, nil
}

func (s *ClubServer) GetClubByID(ctx context.Context, req *clubpb.GetClubByIDRequest) (*clubpb.ClubResponse, error) {
	club, err := s.svc.GetClubByID(ctx, req.IdClub)
	if err != nil {
		return nil, err
	}
	return &clubpb.ClubResponse{
		Club: &clubpb.Club{
			IdClub:   club.IDClub,
			NameClub: club.NameClub,
			CityClub: club.CityClub,
			Budget:   club.Budget,
		},
	}, nil
}

func (s *ClubServer) AddIncome(ctx context.Context, req *clubpb.AddIncomeRequest) (*clubpb.ClubResponse, error) {
	club, err := s.svc.AddIncome(ctx, req.IdClub, req.Amount)
	if err != nil {
		return nil, err
	}
	return &clubpb.ClubResponse{
		Club: &clubpb.Club{
			IdClub:   club.IDClub,
			NameClub: club.NameClub,
			CityClub: club.CityClub,
			Budget:   club.Budget,
		},
	}, nil
}

func (s *ClubServer) Spend(ctx context.Context, req *clubpb.SpendRequest) (*clubpb.ClubResponse, error) {
	club, err := s.svc.Spend(ctx, req.IdClub, req.Amount)
	if err != nil {
		return nil, err
	}
	return &clubpb.ClubResponse{
		Club: &clubpb.Club{
			IdClub:   club.IDClub,
			NameClub: club.NameClub,
			CityClub: club.CityClub,
			Budget:   club.Budget,
		},
	}, nil
}
