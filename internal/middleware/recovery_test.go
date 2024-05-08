package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

func TestRecovery(t *testing.T) {
	// Gestionnaire qui déclenche une `panic`
	panicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("Intentional panic for testing")
	})

	// Créer le middleware de récupération autour du gestionnaire
	handler := middleware.Recovery(panicHandler)

	// Créer un enregistreur de réponse HTTP et une requête simulée
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Effectuer la requête
	handler.ServeHTTP(rr, req)

	// Vérifier si la réponse est un statut 500
	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, got %d", status)
	}

	// Optionnel: Vérifiez le contenu de la réponse
	expectedBody := "Internal Server Error"
	if !strings.Contains(rr.Body.String(), expectedBody) {
		t.Errorf("Expected body to contain %q, got %q", expectedBody, rr.Body.String())
	}
}
