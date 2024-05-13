package server

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

// InitTemplate initialise les templates et configure le serveur de fichiers statiques
var tpl *template.Template

func InitTemplates() {
	// Charger les templates absolus
	absTemplatesPath, err := filepath.Abs("view/templates/**/*.html")
	if err != nil {
		log.Fatalf("Error getting absolute path for templates: %v", err)
	}
	tpl = template.Must(template.ParseGlob(absTemplatesPath))

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
	InitTemplates()

	protectedPaths := []string{"/admin", "/user", "/post-sentiment", "/sentiment-comment", "/create-post-submit"}

	// Créez un multiplexer
	mux := http.NewServeMux()
	mux.Handle("/view/assets/", http.StripPrefix("/view/assets/", http.FileServer(http.Dir("view/assets"))))

	mux.Handle("/", controller.Home(db, tpl))
	mux.Handle("/test", controller.Test(db, tpl))
	mux.Handle("/filtered-home", controller.FilteredHome(db, tpl))
	mux.Handle("/postbyid", controller.UniquePost(db, tpl))
	mux.Handle("/submit-signup", controller.SignUpSubmission(db, tpl))
	mux.Handle("/signup", controller.SignUp(db, tpl))
	mux.Handle("/signin", controller.SignIn(tpl))
	mux.Handle("/submit-signin", controller.SignInSubmit(db, tpl))
	mux.Handle("/disconnect", controller.Disconnect(db, tpl))
	mux.Handle("/post-sentiment", controller.SentimentPost(db, tpl))
	mux.Handle("/comment-sentiment", controller.SentimentComment(db, tpl))
	mux.Handle("/filtered-sentiments", controller.FilterByReact(db, tpl))
	mux.Handle("/create-post", controller.NewPost(db, tpl))
	mux.Handle("/create-post-submit", controller.SubmitPostHandler(db, tpl))

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
