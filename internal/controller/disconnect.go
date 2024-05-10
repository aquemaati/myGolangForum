package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/aquemaati/myGolangForum.git/internal/model"
)

func Disconnect(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		sessionID, ok := r.Context().Value(middleware.SessionIdContextKey).(string)
		if !ok {
			// Handle the case where the userID is missing or of an incorrect type
			log.Println("Session ID not found in context", http.StatusUnauthorized)
		}

		fmt.Println("session id is (in disconect)", sessionID)
		model.InvalidateSessionCookie(w)
		model.DeleteSession(db, sessionID)

		// ne pas oublier de clear le cache
		tpl.ExecuteTemplate(w, "index.html", nil)
	}
}
