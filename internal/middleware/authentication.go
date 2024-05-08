package middleware

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"sync"
	"time"
)

type contextKey string

const userContextKey contextKey = "userID"

type SessionCache struct {
	sync.Map
	expiration time.Duration
}

func NewSessionCache(expiration time.Duration) *SessionCache {
	return &SessionCache{expiration: expiration}
}

// Set ajoute une session dans le cache
func (sc *SessionCache) Set(token string, userID int) {
	sc.Store(token, userID)
	// Expire l'entrée après la durée spécifiée
	go func() {
		time.Sleep(sc.expiration)
		sc.Delete(token)
	}()
}

// Get récupère une session depuis le cache
func (sc *SessionCache) Get(token string) (int, bool) {
	if userID, found := sc.Load(token); found {
		return userID.(int), true
	}
	return 0, false
}

func AuthenticationWithCache(db *sql.DB, cache *SessionCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Récupérer le cookie contenant le jeton d'authentification
			cookie, err := r.Cookie("session_token")
			if err != nil {
				if err == http.ErrNoCookie {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
				} else {
					log.Printf("Error retrieving session cookie: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Récupérer la valeur du cookie (le jeton)
			token := cookie.Value

			// Vérifier si le token existe dans le cache
			if userID, found := cache.Get(token); found {
				// Ajouter l'ID utilisateur au contexte
				ctx := context.WithValue(r.Context(), userContextKey, userID)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Si le token n'est pas dans le cache, interroger la base de données
			var userID int
			err = db.QueryRow("SELECT user_id FROM sessions WHERE token = ?", token).Scan(&userID)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
				} else {
					log.Printf("Error querying database for session token: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			// Ajouter le token et l'utilisateur au cache
			cache.Set(token, userID)

			// Ajouter l'ID utilisateur au contexte
			ctx := context.WithValue(r.Context(), userContextKey, userID)
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
