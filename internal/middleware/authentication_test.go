package middleware_test

import (
	"database/sql"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func TestAuthentication(t *testing.T) {
	// Créer une base de données SQLite en mémoire
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Unable to open database: %v", err)
	}
	defer db.Close()

	// Créer une table de sessions de test
	_, err = db.Exec(`CREATE TABLE Sessions (sessionId TEXT, userId TEXT)`)
	if err != nil {
		t.Fatalf("Unable to create Sessions table: %v", err)
	}

	// Insérer un jeton de session de test
	_, err = db.Exec(`INSERT INTO Sessions (sessionId, userId) VALUES (?, ?)`, "abc123", "user1")
	if err != nil {
		t.Fatalf("Unable to insert test session: %v", err)
	}

	// Créer un cache de sessions
	sessionCache := middleware.NewSessionCache(10 * time.Minute)

	// Initialiser les middlewares
	cacheHandler := middleware.CacheHandler(sessionCache)
	authMiddleware := middleware.Authentication(db, sessionCache)

	// Créer un gestionnaire final qui vérifie le contexte de l'utilisateur
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIdContextKey).(string)
		if !ok {
			log.Println("No user ID found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		w.Write([]byte("User ID: " + userID))
	})

	// Envelopper le gestionnaire final avec les middlewares
	handler := cacheHandler(authMiddleware(finalHandler))

	// Simuler une requête HTTP GET avec le cookie de session
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{
		Name:  "session_token",
		Value: "abc123",
	})

	// Créer un enregistreur de réponse HTTP
	rr := httptest.NewRecorder()

	// Effectuer la requête via le gestionnaire
	handler.ServeHTTP(rr, req)

	// Vérifier le statut de la réponse (200 OK attendu)
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rr.Code)
	}

	// Vérifier que le corps de la réponse contient l'ID utilisateur
	expected := "User ID: user1"
	if rr.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, rr.Body.String())
	}
}
