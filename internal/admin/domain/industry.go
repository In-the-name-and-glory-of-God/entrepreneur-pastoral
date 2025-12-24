package domain

// Industry corresponds to the "industries" table.
// The Key field contains a translation key (e.g., "industry.technology").
type Industry struct {
	ID  int16  `json:"id" db:"id"`
	Key string `json:"key" db:"key"`
}
