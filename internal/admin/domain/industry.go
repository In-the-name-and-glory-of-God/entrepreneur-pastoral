package domain

// Industry corresponds to the "industries" table.
type Industry struct {
	ID   int16  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
