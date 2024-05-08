package server

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/aquemaati/myGolangForum.git/internal/config"
	"github.com/aquemaati/myGolangForum.git/internal/controller"
	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

func InitializeServer(envFilePath string) (*http.Server, error) {
	// Load environment variables
	if err := config.LoadEnvFile(envFilePath); err != nil {
		return nil, err
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", controller.Home)

	handler := middleware.Recovery(middleware.SecurityHeaders(middleware.Logging(mux)))

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
