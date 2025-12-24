package domain

// FieldOfWork corresponds to the "fields_of_work" table.
// The Key field contains a translation key (e.g., "field_of_work.technology").
type FieldOfWork struct {
	ID  int16  `json:"id" db:"id"`
	Key string `json:"key" db:"key"`
}
