package repository

import (
	"Transfer/entities"
	"context"
	"database/sql"
)

type TransferRepository interface {
	Create(ctx context.Context, t *entities.Transfer) error
	GetByID(ctx context.Context, id int64) (*entities.Transfer, error)
	GetAll(ctx context.Context) ([]*entities.Transfer, error)
	Update(ctx context.Context, t *entities.Transfer) error
	Delete(ctx context.Context, id int64) error
}

type transferRepository struct {
	db *sql.DB
}

func NewTransferRepository(db *sql.DB) TransferRepository {
	return &transferRepository{db: db}
}

func (r *transferRepository) Create(ctx context.Context, t *entities.Transfer) error {
	query := `
        INSERT INTO transfer (player_id, old_club_id, new_club_id, transfer_fee, transfer_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id_transfer
    `
	return r.db.QueryRowContext(ctx, query,
		t.PlayerID, t.OldClubID, t.NewClubID, t.TransferFee, t.TransferDate,
	).Scan(&t.IDTransfer)
}

func (r *transferRepository) GetByID(ctx context.Context, id int64) (*entities.Transfer, error) {
	var t entities.Transfer
	query := `
        SELECT id_transfer, player_id, old_club_id, new_club_id, transfer_fee, transfer_date
        FROM transfer WHERE id_transfer = $1
    `
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.IDTransfer, &t.PlayerID, &t.OldClubID, &t.NewClubID, &t.TransferFee, &t.TransferDate,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &t, err
}

func (r *transferRepository) GetAll(ctx context.Context) ([]*entities.Transfer, error) {
	query := `
        SELECT id_transfer, player_id, old_club_id, new_club_id, transfer_fee, transfer_date
        FROM transfer ORDER BY id_transfer ASC
    `
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transfers []*entities.Transfer
	for rows.Next() {
		var t entities.Transfer
		if err := rows.Scan(
			&t.IDTransfer, &t.PlayerID, &t.OldClubID, &t.NewClubID, &t.TransferFee, &t.TransferDate,
		); err != nil {
			return nil, err
		}
		transfers = append(transfers, &t)
	}
	return transfers, nil
}

func (r *transferRepository) Update(ctx context.Context, t *entities.Transfer) error {
	query := `
        UPDATE transfer
        SET player_id=$1, old_club_id=$2, new_club_id=$3, transfer_fee=$4, transfer_date=$5
        WHERE id_transfer=$6
    `
	_, err := r.db.ExecContext(ctx, query,
		t.PlayerID, t.OldClubID, t.NewClubID, t.TransferFee, t.TransferDate, t.IDTransfer,
	)
	return err
}

func (r *transferRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM transfer WHERE id_transfer=$1`, id)
	return err
}
