package main

import (
	"log"
	"os"

	"github.com/aquemaati/myGolangForum.git/internal/server"
)

func main() {
	server, err := server.InitializeServer(".env", "database/database.db")
	if err != nil {
		log.Fatalf("error initializing server: %v\n", err)
		return
	}

	certFile := os.Getenv("TLS_CERT")
	keyFile := os.Getenv("TLS_KEY")
	if certFile != "" && keyFile != "" {
		log.Printf("Server listening on https://localhost%s", server.Addr)
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			log.Fatalf("Error starting TLS server: %v", err)
		}
	} else {
		log.Printf("Server listening on http://localhost%s", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}
}
