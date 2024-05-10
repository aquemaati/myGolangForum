package controller

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

func SignIn(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tpl.ExecuteTemplate(w, "signin.html", nil)
	}
}

func SignInSubmit(tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)
		fmt.Println(formData)
	}
}
