package service

import (
	"Club/entities"
	"Club/repository"
	"context"
	"errors"
)

type ClubService interface {
	CreateClub(ctx context.Context, c *entities.Club) error
	GetClubByID(ctx context.Context, id int64) (*entities.Club, error)
	GetAllClubs(ctx context.Context) ([]*entities.Club, error)
	UpdateClub(ctx context.Context, c *entities.Club) error
	DeleteClub(ctx context.Context, id int64) error

	// aturan budget
	AddIncome(ctx context.Context, id int64, amount int64) (*entities.Club, error)
	Spend(ctx context.Context, id int64, amount int64) (*entities.Club, error)
}

type clubService struct {
	repo repository.ClubRepository
}

func NewClubService(repo repository.ClubRepository) ClubService {
	return &clubService{repo: repo}
}

// ================= CRUD =================

func (s *clubService) CreateClub(ctx context.Context, c *entities.Club) error {
	if c.Budget < 0 {
		return errors.New("budget awal tidak boleh negatif")
	}
	return s.repo.Create(ctx, c)
}

func (s *clubService) GetClubByID(ctx context.Context, id int64) (*entities.Club, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *clubService) GetAllClubs(ctx context.Context) ([]*entities.Club, error) {
	return s.repo.GetAll(ctx)
}

func (s *clubService) UpdateClub(ctx context.Context, c *entities.Club) error {
	// budget tidak bisa diedit langsung lewat update
	old, err := s.repo.GetByID(ctx, c.IDClub)
	if err != nil {
		return err
	}
	if old == nil {
		return errors.New("club tidak ditemukan")
	}
	c.Budget = old.Budget // pastikan budget tidak berubah
	return s.repo.Update(ctx, c)
}

func (s *clubService) DeleteClub(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

// ================= BUDGET =================

func (s *clubService) AddIncome(ctx context.Context, id int64, amount int64) (*entities.Club, error) {
	club, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	club.Budget += amount
	if err := s.repo.Update(ctx, club); err != nil {
		return nil, err
	}
	return club, nil
}

func (s *clubService) Spend(ctx context.Context, id int64, amount int64) (*entities.Club, error) {
	club, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if club.Budget < amount {
		return nil, errors.New("budget tidak cukup")
	}
	club.Budget -= amount
	if err := s.repo.Update(ctx, club); err != nil {
		return nil, err
	}
	return club, nil
}
