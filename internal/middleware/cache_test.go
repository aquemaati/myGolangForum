package middleware_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

func TestCacheHandler(t *testing.T) {
	// Créer un cache de sessions
	sessionCache := middleware.NewSessionCache(10 * time.Minute)

	// Ajouter un jeton et un ID utilisateur au cache
	sessionCache.Set("abc123", "user1")

	// Initialiser le CacheHandler
	cacheHandler := middleware.CacheHandler(sessionCache)

	// Créer un gestionnaire final qui vérifie le contexte
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vérifier que l'ID utilisateur et le token sont dans le contexte
		userID, ok := r.Context().Value(middleware.UserIdContextKey).(string)
		if !ok {
			log.Println("No user ID found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionID, ok := r.Context().Value(middleware.SessionIdContextKey).(string)
		if !ok {
			log.Println("No session ID found in context")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Vérifier que les valeurs correspondent
		if userID != "user1" || sessionID != "abc123" {
			http.Error(w, "Mismatched user ID or session ID", http.StatusUnauthorized)
			return
		}

		w.Write([]byte("Cache hit, user ID: " + userID))
	})

	// Envelopper le gestionnaire final avec le CacheHandler
	handler := cacheHandler(finalHandler)

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
	expected := "Cache hit, user ID: user1"
	if rr.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, rr.Body.String())
	}
}
