package middleware

import (
	"log"
	"net/http"
)

// Recovery handles critical errors (panics).
//
// Prevents the server from crashing by recovering from any panics encountered.
// Returns a 500 error to the client.
// Captures panics throughout the middleware chain as well as those in the final handlers.
// If panic captured, all process stop and 500 is sent to the client
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Return a 500 error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				// Log the captured error for debugging purposes
				log.Printf("Recovered from panic: %v", err)
			}
		}()
		// Proceed with the next handler in the chain if no panic is captured
		next.ServeHTTP(w, r)
	})
}

/*
Fonctionnement du Middleware de Récupération :
Encapsulation:
RecoveryMiddleware encapsule toute la chaîne de gestionnaires, y compris d'autres middlewares et les gestionnaires finaux.
C'est le seul middleware à avoir le bloc defer pour capturer les panics de manière centralisée.
Rôle du defer et de recover:
Le bloc différé (defer) s'assure que recover est exécuté après chaque requête traitée.
Si une panic se produit à n'importe quel point de la chaîne de gestionnaires, recover la capture et renvoie une réponse HTTP 500.
Traitement Normal:
Si aucune panic n'est capturée, la requête continue son exécution normale via next.ServeHTTP.

*/
