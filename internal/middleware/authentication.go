// middleware/authentication.go
package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
)

func Authentication(db *sql.DB, cache *SessionCache, protectedPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Récupérer le chemin d'URL de la requête
			requestPath := r.URL.Path

			// Vérifier si le chemin d'URL nécessite une authentification
			requiresAuthentication := false
			for _, path := range protectedPaths {
				if strings.HasPrefix(requestPath, path) {
					requiresAuthentication = true
					break
				}
			}

			// Si le chemin d'URL nécessite une authentification
			if requiresAuthentication {
				// Récupérer le token de session du contexte
				token, ok := r.Context().Value(SessionIdContextKey).(string)
				if !ok {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				// Vérifier si l'ID utilisateur est déjà dans le contexte
				if _, ok := r.Context().Value(UserIdContextKey).(string); ok {
					// Si trouvé dans le contexte, continuer
					next.ServeHTTP(w, r)
					return
				}

				// Interroger la base de données pour obtenir l'ID utilisateur
				var userID string
				err := db.QueryRow("SELECT userId FROM Sessions WHERE sessionId = ?", token).Scan(&userID)
				if err != nil {
					if err == sql.ErrNoRows {
						http.Error(w, "Unauthorized", http.StatusUnauthorized)
					} else {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					}
					return
				}

				// Ajouter l'ID utilisateur au contexte
				ctx := context.WithValue(r.Context(), UserIdContextKey, userID)

				// Mettre à jour le cache avec le nouveau token et l'ID utilisateur
				cache.Set(token, userID)

				// Continuer avec le gestionnaire suivant
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Si le chemin d'URL ne nécessite pas d'authentification, passer à la requête suivante sans effectuer d'authentification
			next.ServeHTTP(w, r)
		})
	}
}

/*
Oui, un client peut avoir besoin d'accéder simultanément à de nombreuses données mises en cache. C'est précisément pourquoi l'utilisation d'une structure concurrente comme `sync.Map` ou d'un autre mécanisme de cache est utile. Voici les raisons pour lesquelles un cache efficace est essentiel dans ce contexte :

### Besoins du Client pour un Cache Simultané :

1. **Accès Concurentiel:**
   - Plusieurs goroutines (ou requêtes) du client peuvent tenter d'accéder ou de mettre à jour différentes données en même temps.
   - Le cache doit donc être capable de gérer ces accès concurrents sans créer de conflits ou ralentir les performances.

2. **Récupération Rapide:**
   - Le cache doit être capable de fournir rapidement les données nécessaires pour réduire la latence des requêtes.
   - C'est essentiel pour des tâches comme l'authentification, les autorisations, ou le chargement de configurations.

3. **Réduction de la Charge sur la Base de Données:**
   - Un cache bien géré réduit considérablement le nombre d'appels directs à la base de données, ce qui améliore les performances globales.

### Stratégies pour un Cache Simultané :

1. **`sync.Map`:**
   - Pratique pour les scénarios simples où des paires clé-valeur peuvent être stockées en parallèle.

2. **Mutex et RWMutex:**
   - Pour des structures plus complexes (par exemple, des maps imbriquées), des verrous comme `sync.Mutex` ou `sync.RWMutex` peuvent être utilisés pour assurer une synchronisation fine.

3. **Caches Tiers:**
   - Des solutions comme Redis ou Memcached permettent d'avoir un cache distribué accessible depuis plusieurs services ou instances.

### Exemples de Scénarios Simultanés :

1. **Session et Authentification:**
   - Vérification des sessions utilisateur lors de multiples connexions simultanées.

2. **Profil Utilisateur:**
   - Chargement de profils utilisateur et de préférences pour les différents clients.

3. **Données de Configuration:**
   - Fourniture rapide des configurations partagées entre différentes requêtes.

### Résumé :

- Les clients ont souvent besoin d'accéder simultanément à diverses informations.
- Un cache bien conçu peut gérer ces demandes en fournissant un accès rapide et concurrentiel.
- Assurez-vous que votre cache peut répondre aux besoins spécifiques du client tout en restant simple et performant.
*/
