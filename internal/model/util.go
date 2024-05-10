package model

import (
	"database/sql"
	"log"
)

// ExecuteQuery est une fonction générique pour exécuter une requête SQL et scanner les résultats dans des structures
func ExecuteQuery[T any](db *sql.DB, query string, scanFunc func(*sql.Rows) (T, error), args ...any) ([]T, error) {
	// Préparer la requête
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing SQL query: %v", err)
		return nil, err
	}
	defer stmt.Close()

	// Exécuter la requête
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scanner les résultats
	var results []T
	for rows.Next() {
		record, err := scanFunc(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, record)
	}

	// Vérifier les erreurs rencontrées pendant l'itération
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

//In summary, use db.QueryRow for queries that you know will always return one row or expect a unique
