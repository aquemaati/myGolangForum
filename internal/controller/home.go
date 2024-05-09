package controller

import (
	"database/sql"
	"net/http"
)

// HomeHandler handles the root path
func Home(db *sql.DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implémentation de la page d'accueil
		// Exemple simple d'une réponse HTML
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<h1>Welcome to the Home Page!</h1>"))
	})
}
