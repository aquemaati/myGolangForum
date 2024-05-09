package server

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/aquemaati/myGolangForum.git/database"
	"github.com/aquemaati/myGolangForum.git/internal/config"
	"github.com/aquemaati/myGolangForum.git/internal/controller"
	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

// Initialise le cache de session
func initSessionCache() *middleware.SessionCache {
	return middleware.NewSessionCache(10 * time.Minute)
}

// Initialise le serveur HTTP
func InitializeServer(envFilePath, dbPath string) (*http.Server, error) {
	// Chargez les variables d'environnement
	if err := config.LoadEnvFile(envFilePath); err != nil {
		return nil, err
	}

	// Récupérez le port depuis l'environnement
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialise la base de données et le cache
	db, err := database.InitDatabase(dbPath)
	if err != nil {
		return nil, err
	}
	sessionCache := initSessionCache()

	protectedPaths := []string{"/admin", "/user"}

	// Créez un multiplexer
	mux := http.NewServeMux()
	mux.Handle("/", controller.Home(db))

	// Chaîne de middlewares
	handler := middleware.Recovery(
		middleware.SecurityHeaders(
			middleware.Logging(
				middleware.ParseForm( // Ajoutez le parsing du formulaire avant les autres middlewares
					middleware.CacheHandler(sessionCache)(
						middleware.Authentication(db, sessionCache, protectedPaths)(mux),
					),
				),
			),
		),
	)

	// Configurez le serveur avec les paramètres TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	server := &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      tlsConfig,
	}

	return server, nil
}
