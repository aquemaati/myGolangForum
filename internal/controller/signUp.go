package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
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

		fmt.Println(username, email, password, confirmPassword)

		tpl.ExecuteTemplate(w, "index.html", nil)
	}
}
