package model

import (
	"database/sql"
	"log"
	"net/http"
	"time"
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

// ExecuteNonQuery is a generic function to execute non-SELECT SQL statements
func ExecuteNonQuery(db *sql.DB, query string, args ...any) (int64, error) {
	// Prepare the query
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing SQL query: %v", err)
		return 0, err
	}
	defer stmt.Close()

	// Execute the query
	result, err := stmt.Exec(args...)
	if err != nil {
		log.Printf("Error executing SQL query: %v", err)
		return 0, err
	}

	// Return the number of rows affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

// SetCookie sets a cookie with common session attributes
func SetCookie(w http.ResponseWriter, name, value string, expires time.Time, secure bool, path string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Expires:  expires,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode, // Adjust as necessary
	}
	http.SetCookie(w, cookie)
}

// ExecuteSingleQuery is a generic function to execute a SQL query and scan a single row into a structure
func ExecuteSingleQuery[T any](db *sql.DB, query string, scanFunc func(*sql.Row) (T, error), args ...any) (T, error) {
	var result T

	// Prepare the statement
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("Error preparing SQL query: %v", err)
		return result, err
	}
	defer stmt.Close()

	// Execute the query and scan the result
	row := stmt.QueryRow(args...)
	result, err = scanFunc(row)
	if err != nil {
		log.Printf("Error scanning SQL query result: %v", err)
		return result, err
	}

	return result, nil
}
