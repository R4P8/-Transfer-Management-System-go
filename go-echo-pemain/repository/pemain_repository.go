package repository

import (
	"context"
	"database/sql"
	"errors"

	"Pemain/entities"
)

type PemainRepository interface {
	Create(ctx context.Context, p *entities.Pemain) error
	GetByID(ctx context.Context, id int64) (*entities.Pemain, error)
	GetAll(ctx context.Context) ([]*entities.Pemain, error)
	UpdateByPlayer(ctx context.Context, id int64, status string) error
	UpdateByClub(ctx context.Context, id int64, clubID int64, marketValue float64, status string) error
	Delete(ctx context.Context, id int64) error
}

type pemainRepository struct {
	db *sql.DB
}

func NewPemainRepository(db *sql.DB) PemainRepository {
	return &pemainRepository{db: db}
}

// Create pemain baru
func (r *pemainRepository) Create(ctx context.Context, p *entities.Pemain) error {

	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM pemain WHERE id_pemain=$1)`
	err := r.db.QueryRowContext(ctx, checkQuery, checkQuery, p.IDPemain)
	if err == nil && exists {
		return errors.New("id_pemain already exists")
	}

	query := `
		INSERT INTO pemain (name, club_id, market_value, status_pemain)
		VALUES ($1, $2, $3, $4)
		RETURNING id_pemain
	`

	// default value kalau kosong
	if p.ClubID == nil {
		p.ClubID = nil // free agent
	}
	if p.MarketValue == 0 {
		p.MarketValue = 0
	}
	if p.Status == "" {
		p.Status = "free_agent"
	}

	return r.db.QueryRowContext(ctx, query, p.Name, p.ClubID, p.MarketValue, p.Status).
		Scan(&p.IDPemain)
}

// Ambil pemain berdasarkan ID
func (r *pemainRepository) GetByID(ctx context.Context, id int64) (*entities.Pemain, error) {
	var p entities.Pemain
	var clubID sql.NullInt64

	query := `SELECT id_pemain, name, club_id, market_value, status_pemain 
	          FROM pemain WHERE id_pemain=$1`

	// QueryRowContext → Scan → baru ada error
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.IDPemain,
		&p.Name,
		&clubID,
		&p.MarketValue,
		&p.Status,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if clubID.Valid {
		p.ClubID = &clubID.Int64
	} else {
		p.ClubID = nil
	}

	return &p, nil
}

// Ambil semua pemain
func (r *pemainRepository) GetAll(ctx context.Context) ([]*entities.Pemain, error) {
	var players []*entities.Pemain
	query := `SELECT id_pemain, name, club_id, market_value, status_pemain FROM pemain ORDER BY id_pemain ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return players, nil
}

// Update oleh pemain sendiri (hanya status)
func (r *pemainRepository) UpdateByPlayer(ctx context.Context, id int64, status string) error {
	query := `UPDATE pemain SET status_pemain=$1 WHERE id_pemain=$2`
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}

// Update oleh club (boleh club_id, market_value, status)
func (r *pemainRepository) UpdateByClub(ctx context.Context, id int64, clubID int64, marketValue float64, status string) error {
	query := `UPDATE pemain SET club_id=$1, market_value=$2, status_pemain=$3 WHERE id_pemain=$4`
	_, err := r.db.ExecContext(ctx, query, clubID, marketValue, status, id)
	return err
}

// Delete pemain
func (r *pemainRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM pemain WHERE id_pemain=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
