package middleware

import (
	"context"
	"log"
	"net/http"
)

func ParseForm(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var formData map[string][]string

		switch r.Method {
		case http.MethodGet:
			formData = r.URL.Query()

		case http.MethodPost:
			contentType := r.Header.Get("Content-Type")
			if contentType == "" {
				log.Println("No Content-Type header")
				http.Error(w, "Content-Type header required", http.StatusBadRequest)
				return
			}

			if contentType == "application/x-www-form-urlencoded" {
				if err := r.ParseForm(); err != nil {
					log.Printf("Error parsing form: %v", err)
					http.Error(w, "Unable to parse form", http.StatusBadRequest)
					return
				}
				formData = r.Form

			} else if contentType == "multipart/form-data" || contentType[:19] == "multipart/form-data" {
				if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB de taille maximale
					log.Printf("Error parsing multipart form: %v", err)
					http.Error(w, "Unable to parse multipart form", http.StatusBadRequest)
					return
				}
				formData = r.MultipartForm.Value

			} else {
				log.Printf("Unsupported content type: %s", contentType)
				http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
				return
			}

		default:
			log.Println("Method not allowed")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Ajouter les données du formulaire au contexte
		ctx := context.WithValue(r.Context(), FormDataKey, formData)

		// Continuer avec le prochain gestionnaire
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// in the handler, use formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)
/*
/ Utiliser les données du formulaire
	name := formData["name"][0]
	email := formData["email"][0]
*/
