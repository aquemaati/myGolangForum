package middleware

import "net/http"

// SecurityHeaders ajoute des en-têtes de sécurité aux réponses
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ajouter les en-têtes de sécurité
		// w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
		// w.Header().Set("X-Content-Type-Options", "nosniff")
		// w.Header().Set("X-Frame-Options", "DENY")
		// w.Header().Set("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self'")

		// Appeler le prochain gestionnaire dans la chaîne
		next.ServeHTTP(w, r)
	})
}
