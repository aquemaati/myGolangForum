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

		// create session with JWT
		jwt, err := model.GenerateJWT(db, user.ID)
		if err != nil {
			log.Println("could not create session in database", err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
		// create cookie
		model.SetCookie(w, "session_token", jwt, time.Now().Add(24*time.Hour), true, "/")
		// Redirect to the last accessed page or default to a home/dashboard page
		redirectURL := r.URL.Query().Get("redirect")
		if redirectURL == "" {
			redirectURL = "/" // Default page if no specific redirect provided
		}
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
