package database

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// Initialise la base de données (fonction simplifiée)
func InitDatabase(dbPath string) (*sql.DB, error) {
	// Obtenir le chemin absolu
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, err
	}

	// Vérifier si le fichier de la base de données existe
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, errors.New("database file does not exist")
	}

	// Ouvrir la connexion à la base de données
	db, err := sql.Open("sqlite3", absPath)
	if err != nil {
		return nil, err
	}

	// Tester la connexion
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
