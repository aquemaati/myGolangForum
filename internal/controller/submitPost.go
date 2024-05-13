package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

func SubmitPoast(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)
		fmt.Println(formData)
	}
}
