package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/aquemaati/myGolangForum.git/internal/model"
)

func SignIn(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "signin.html", nil)
	}
}

func SignInSubmit(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)

		//Step for signin
		// check credentials (username, password)

		email := formData["email"][0]
		password := formData["password"][0]
		fmt.Println("look here")
		fmt.Println(email, password)

		user, err := model.CheckUserSignIn(db, email, password)
		if err != nil {
			http.Error(w, "error whilee connecting", http.StatusBadRequest)
			return
		}
		// create session
		sess, err := model.CreateSession(db, user.ID, 24*time.Hour)
		if err != nil {
			log.Println("could not create session in database", err)
		}
		// create cookie
		model.SetCookie(w, "session_token", sess.SessionID, sess.ExpiresAt, true, "/")
		// Redirect to the last accessed page or default to a home/dashboard page
		redirectURL := r.URL.Query().Get("redirect")
		if redirectURL == "" {
			redirectURL = "/" // Default page if no specific redirect provided
		}
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}
