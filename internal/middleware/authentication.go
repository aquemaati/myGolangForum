// middleware/authentication.go
package middleware

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/aquemaati/myGolangForum.git/internal/model"
)

// Authentication middleware
func Authentication(db *sql.DB, cache *SessionCache, protectedPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// checking if the path require authentication
			requiresAuth := false
			for _, path := range protectedPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					requiresAuth = true
					break
				}
			}

			// Retrieve the JWT from the cookie
			cookie, err := r.Cookie("session_token") // Assuming 'session_token' is the name of the cookie containing the JWT
			if err != nil && requiresAuth {
				// No token provided and path requires authentication
				http.Redirect(w, r, "/signin", http.StatusFound)
				return
			} else if err != nil {
				// No token provided but path does not require authentication
				next.ServeHTTP(w, r)
				return
			}

			token := cookie.Value
			fmt.Println(token)

			// Parse and validate the JWT
			claims, err := model.ParseJWT(token)
			if err != nil && requiresAuth {
				http.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
				return
			} else if err != nil {
				// Invalid token but path does not require authentication
				next.ServeHTTP(w, r)
				return
			}
			fmt.Println(claims)
			userID, ok := claims[string(UserIdContextKey)].(string)
			if !ok && requiresAuth {
				http.Error(w, "Unauthorized: Unable to parse userID", http.StatusUnauthorized)
				return
			} else if !ok {
				// Cannot parse userID but path does not require authentication
				next.ServeHTTP(w, r)
				return
			}
			fmt.Println(userID)

			// Verify if the session is still active
			var sessionCount int
			err = db.QueryRow("SELECT COUNT(*) FROM Sessions WHERE JWT = ? AND UserID = ? AND ExpiresAt > CURRENT_TIMESTAMP", token, userID).Scan(&sessionCount)
			fmt.Println("----session count -----", sessionCount, "-----")
			if err != nil || sessionCount == 0 {
				if requiresAuth {
					http.Error(w, "Unauthorized: Session is not active or does not exist", http.StatusUnauthorized)
				} else {
					fmt.Println("Non-authenticated path accessed with invalid session")
					next.ServeHTTP(w, r)
				}
				return
			}

			fmt.Println("here is the context")
			// Pass user ID to the context of the next request
			ctx := context.WithValue(r.Context(), UserIdContextKey, userID)
			fmt.Println("hello from authen")
			next.ServeHTTP(w, r.WithContext(ctx))
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
