package service

import (
	"context"
	"errors"

	"Pemain/entities"
	"Pemain/repository"
)

type PemainService interface {
	Create(ctx context.Context, p *entities.Pemain) error
	GetByID(ctx context.Context, id int64) (*entities.Pemain, error)
	Update(ctx context.Context, p *entities.Pemain) error
	Delete(ctx context.Context, id int64) error
}

type pemainService struct {
	repo repository.PemainRepository
}

func NewPemainService(repo repository.PemainRepository) PemainService {
	return &pemainService{repo: repo}
}

// Create pemain baru
func (s *pemainService) Create(ctx context.Context, p *entities.Pemain) error {
	if p.Name == "" {
		return errors.New("nama pemain wajib diisi")
	}
	return s.repo.Create(ctx, p)
}

// Get pemain by ID
func (s *pemainService) GetByID(ctx context.Context, id int64) (*entities.Pemain, error) {
	return s.repo.GetByID(ctx, id)
}

// Update pemain (aturan:
// - pemain boleh update status
// - club boleh update club_id & market_value)
func (s *pemainService) Update(ctx context.Context, p *entities.Pemain) error {
	existing, err := s.repo.GetByID(ctx, p.IDPemain)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("pemain tidak ditemukan")
	}

	// Validasi aturan update
	if p.ClubID != nil && *p.ClubID != *existing.ClubID {
		// hanya club yang boleh ganti club_id → misalnya validasi di layer lain
		return errors.New("club_id hanya bisa diubah oleh club")
	}

	if p.MarketValue != existing.MarketValue {
		// hanya club yang boleh ubah market value
		return errors.New("market value hanya bisa diubah oleh club")
	}

	// status pemain boleh diubah sendiri
	return s.repo.UpdateByPlayer(ctx, p)
}

// Delete pemain
func (s *pemainService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
