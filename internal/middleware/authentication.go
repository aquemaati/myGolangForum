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
