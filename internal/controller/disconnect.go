package controller

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/aquemaati/myGolangForum.git/internal/model"
)

func Disconnect(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			log.Println("token not provided")
		}
		// Prepare the SQL statement to avoid SQL injection
		stmt, err := db.Prepare("UPDATE Sessions SET ExpiresAt = ? WHERE JWT = ?")
		if err != nil {
			log.Fatalf("Error preparing update statement: %v", err)
			return
		}
		defer stmt.Close()

		// Execute the statement with specific parameters
		_, err = stmt.Exec(time.Now(), cookie.Value)
		if err != nil {
			log.Printf("Error updating session expiration: %v", err)
		}

		// just sent cookie
		model.SetCookie(w, "session_token", "out", time.Now(), true, "/")

		// invalidate sessions in database en updatant expires at

		tpl.ExecuteTemplate(w, "index.html", nil)
	}
}
