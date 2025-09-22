package entities

type Pemain struct {
	IDPemain    int64   `json:"id_pemain" db:"id_pemain"`
	Name        string  `json:"name" db:"name"`
	ClubID      *int64  `json:"club_id,omitempty" db:"club_id"`
	MarketValue float64 `json:"market_value" db:"market_value"`
	Status      string  `json:"status_pemain" db:"status_pemain"`
}
