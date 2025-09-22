package entities

type Club struct {
	IDClub   int64  `json:"id_club" db:"id_club"`
	NameClub string `json:"name_club" db:"name_club"`
	CityClub string `json:"city_club" db:"city_club"`
	Budget   int64  `json:"budget" db:"budget"`
}
