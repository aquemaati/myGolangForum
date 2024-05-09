package database_test

import (
	"testing"

	"github.com/aquemaati/myGolangForum.git/database"
	_ "github.com/mattn/go-sqlite3"
)

func TestInitDatabase(t *testing.T) {
	dbPath := "database.db"

	db, err := database.InitDatabase(dbPath)
	if err != nil {
		t.Fatalf("Error initializing database: %v", err)
	}
	defer db.Close()

	// la connexion est valide?
	if err := db.Ping(); err != nil {
		t.Errorf("Unable to ping the database: %v", err)
	}
}
