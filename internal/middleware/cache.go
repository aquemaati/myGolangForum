package middleware

import (
	"net/http"
	"sync"
	"time"
)

type SessionCache struct {
	sync.Map
	expiration time.Duration
}

func NewSessionCache(expiration time.Duration) *SessionCache {
	return &SessionCache{expiration: expiration}
}

// Set ajoute une session dans le cache
func (sc *SessionCache) Set(token string, userID string) {
	sc.Store(token, userID)
	// Expire l'entrée après la durée spécifiée
	// supprime automatiquement les sessions out of date
	go func() {
		time.Sleep(sc.expiration)
		sc.Delete(token)
	}()
}

// Get récupère une session depuis le cache
func (sc *SessionCache) Get(token string) (string, bool) {
	if userID, found := sc.Load(token); found {
		return userID.(string), true
	}
	return "", false
}

// this function catch cookie if there is one and check the cache for avoying
// too much sql requests
// CacheHandler ajoute des informations de session au contexte si elles sont dans le cache
func CacheHandler(cache *SessionCache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// // Récupérer le cookie contenant le jeton d'authentification
			// cookie, err := r.Cookie("session_token")
			// if err != nil { // if no cookie
			// 	if err == http.ErrNoCookie {
			// 		log.Println("No session cookie found") // go to next middleware
			// 	} else {
			// 		log.Printf("Error retrieving session cookie: %v", err)
			// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			// 		return
			// 	}
			// } else { //if ther is a cookie
			// 	// Ajouter le jeton de session au contexte
			// 	log.Println("cookie foud")
			// 	ctx := context.WithValue(r.Context(), SessionIdContextKey, cookie.Value)
			// 	// Vérifier si le token existe dans le cache pour obtenur l userId
			// 	if userID, found := cache.Get(cookie.Value); found {
			// 		log.Printf("Session token %s found in cache for user %s", cookie.Value, userID)
			// 		ctx = context.WithValue(ctx, UserIdContextKey, userID)
			// 	} // si pas dans le cache, on passe a authentication pour recuperer lluser id dans
			// 	// la base de donnees

			// 	// Passer le contexte mis à jour
			// 	r = r.WithContext(ctx)
			// }

			// Continuer avec le prochain gestionnaire dans la chaîne
			next.ServeHTTP(w, r)
		})
	}
}

//je n ai peut etre besoin que de l id de session

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
