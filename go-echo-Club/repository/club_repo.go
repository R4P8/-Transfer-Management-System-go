package repository

import (
	"Club/entities"
	"context"
	"database/sql"
)

type ClubRepository interface {
	Create(ctx context.Context, c *entities.Club) error
	GetByID(ctx context.Context, id int64) (*entities.Club, error)
	GetAll(ctx context.Context) ([]*entities.Club, error)
	Update(ctx context.Context, c *entities.Club) error
	Delete(ctx context.Context, id int64) error

	// aturan budget
	AddIncome(ctx context.Context, id int64, amount int64) error
	Spend(ctx context.Context, id int64, amount int64) error
}

type clubRepo struct {
	db *sql.DB
}

func NewClubRepository(db *sql.DB) ClubRepository {
	return &clubRepo{db: db}
}

func (r *clubRepo) Create(ctx context.Context, c *entities.Club) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO club (name_club, city_club, budget) VALUES ($1, $2, $3)`,
		c.NameClub, c.CityClub, c.Budget)
	return err
}

func (r *clubRepo) GetByID(ctx context.Context, id int64) (*entities.Club, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id_club, name_club, city_club, budget FROM club WHERE id_club=$1`, id)

	c := &entities.Club{}
	err := row.Scan(&c.IDClub, &c.NameClub, &c.CityClub, &c.Budget)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return c, err
}

func (r *clubRepo) GetAll(ctx context.Context) ([]*entities.Club, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id_club, name_club, city_club, budget FROM club`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*entities.Club
	for rows.Next() {
		c := &entities.Club{}
		if err := rows.Scan(&c.IDClub, &c.NameClub, &c.CityClub, &c.Budget); err != nil {
			return nil, err
		}
		list = append(list, c)
	}
	return list, nil
}

func (r *clubRepo) Update(ctx context.Context, c *entities.Club) error {
	//  budget TIDAK bisa diupdate sembarangan
	_, err := r.db.ExecContext(ctx,
		`UPDATE club SET name_club=$1, city_club=$2 WHERE id_club=$3`,
		c.NameClub, c.CityClub, c.IDClub)
	return err
}

func (r *clubRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM club WHERE id_club=$1`, id)
	return err
}

// ========== Budget Rules ==========

func (r *clubRepo) AddIncome(ctx context.Context, id int64, amount int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE club SET budget = budget + $1 WHERE id_club=$2`, amount, id)
	return err
}

func (r *clubRepo) Spend(ctx context.Context, id int64, amount int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE club SET budget = budget - $1 WHERE id_club=$2`, amount, id)
	return err
}
