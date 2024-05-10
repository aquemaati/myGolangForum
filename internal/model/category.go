package model

import (
	"database/sql"
)

type Categorie struct {
	Id          int
	Name        string
	Description string
}

func ScanCat(rows *sql.Rows) (Categorie, error) {
	var cat Categorie
	err := rows.Scan(
		&cat.Id,
		&cat.Name,
		&cat.Description,
	)

	if err != nil {
		return Categorie{}, err
	}
	return cat, nil
}

func FetchCat(db *sql.DB) ([]Categorie, error) {
	query := `SELECT id, name, description FROM Categories`
	return ExecuteQuery(db, query, ScanCat)
}
