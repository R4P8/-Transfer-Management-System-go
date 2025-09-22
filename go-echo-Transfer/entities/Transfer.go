package entities

type Transfer struct {
	IDTransfer   int64   `json:"id_transfer" db:"id_transfer"`
	PlayerID     int64   `json:"player_id" db:"player_id"`
	OldClubID    int64   `json:"old_club_id,omitempty" db:"old_club_id"`
	NewClubID    int64   `json:"new_club_id,omitempty" db:"new_club_id"`
	TransferFee  float64 `json:"transfer_fee" db:"transfer_fee"`
	TransferDate string  `json:"transfer_date" db:"transfer_date"`
	Status       string  `db:"status" json:"status"`
}
