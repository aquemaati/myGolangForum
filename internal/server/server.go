package server

import (
	"crypto/tls"
	"database/sql"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/aquemaati/myGolangForum.git/internal/config"
	"github.com/aquemaati/myGolangForum.git/internal/controller"
	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

// Initialise la base de données (fonction simplifiée)
func initDatabase(dbPath string) (*sql.DB, error) {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, err
	}
	return sql.Open("mysql", absPath)
}

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
	db, err := initDatabase(dbPath)
	if err != nil {
		return nil, err
	}
	sessionCache := initSessionCache()

	// Créez un multiplexer
	mux := http.NewServeMux()
	mux.HandleFunc("/", controller.Home)

	// Ajoutez les middlewares
	handler := middleware.Recovery(
		middleware.SecurityHeaders(
			middleware.Logging(
				middleware.AuthenticationWithCache(db, sessionCache)(mux),
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
