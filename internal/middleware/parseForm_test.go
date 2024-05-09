package middleware_test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

func TestParseFormMiddleware(t *testing.T) {
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		formData, ok := r.Context().Value(middleware.FormDataKey).(map[string][]string)
		if !ok {
			t.Log("Form data not found in context")
			http.Error(w, "Form data not found in context", http.StatusInternalServerError)
			return
		}

		// Vérifier les valeurs récupérées
		name := formData["name"][0]
		email := formData["email"][0]

		if name != "Alice" || email != "alice@example.com" {
			t.Log("Incorrect form data")
			http.Error(w, "Incorrect form data", http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	parseFormMiddleware := middleware.ParseForm(finalHandler)

	// Tester le cas POST avec `multipart/form-data`
	t.Run("POST multipart/form-data", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)

		_ = writer.WriteField("name", "Alice")
		_ = writer.WriteField("email", "alice@example.com")
		writer.Close()

		req := httptest.NewRequest(http.MethodPost, "/", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		rr := httptest.NewRecorder()
		parseFormMiddleware.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Logf("Expected status code 200, got %d", rr.Code)
		}
	})
}
