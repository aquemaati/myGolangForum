package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

func TestSecurityHeaders(t *testing.T) {
	// Créez un gestionnaire basique
	helloHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World"))
	})

	// Encapsulez le gestionnaire avec le middleware de sécurité
	handler := middleware.SecurityHeaders(helloHandler)

	// Simulez une requête HTTP et capturez la réponse
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	// Effectuez la requête via le gestionnaire
	handler.ServeHTTP(rr, req)

	// Vérifiez les en-têtes de sécurité
	expectedHeaders := map[string]string{
		"Strict-Transport-Security": "max-age=63072000; includeSubDomains",
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"Content-Security-Policy":   "default-src 'self'; script-src 'self'; style-src 'self'",
	}

	for header, expectedValue := range expectedHeaders {
		actualValue := rr.Header().Get(header)
		if actualValue != expectedValue {
			t.Errorf("Expected %s header to be %q, got %q", header, expectedValue, actualValue)
		}
	}

	// Vérifiez le statut de la réponse
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", status)
	}
}
