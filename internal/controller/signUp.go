package controller

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/aquemaati/myGolangForum.git/internal/model"
)

// rediredirect toward signup page
func SignUp(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "signup.html", nil)
	}
}

func SignUpSubmission(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)

		username := formData["username"][0]
		email := formData["email"][0]
		password := formData["password"][0]
		confirmPassword := formData["confirmPassword"][0]

		if password != confirmPassword {
			http.Error(w, "Password not confirmed", http.StatusBadRequest)
			return
		}

		user, err := model.CreateUserInDB(db, username, email, password)
		if err != nil {
			http.Error(w, "could Not add an user"+err.Error(), http.StatusBadRequest)
			return
		}

		//creating a session
		sess, err := model.CreateSession(db, user.ID, 24*time.Hour)
		if err != nil {
			http.Error(w, "Could not insert new session in database: "+err.Error(), http.StatusBadRequest)
			return
		}

		//set cookie
		model.SetCookie(w, "session_token", sess.SessionID, sess.ExpiresAt, true, "/")

		log.Println("new user created: ", username)
		tpl.ExecuteTemplate(w, "index.html", nil)
	}
}
