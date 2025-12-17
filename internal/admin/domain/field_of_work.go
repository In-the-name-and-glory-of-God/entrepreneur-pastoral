package domain

// FieldOfWork corresponds to the "fields_of_work" table.
type FieldOfWork struct {
	ID   int16  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
