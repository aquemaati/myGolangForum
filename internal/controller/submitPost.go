package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
)

func SubmitPoast(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		index, err := UserAuthParser(r, db, w)
		if err != nil {
			log.Println(err)
		}

		userId := index.UserID

		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)
		fmt.Println(formData)

		titlePost := formData["titlePost"][0]
		descriptionPost := formData["descriptionPost"][0]
		contentPost := formData["contentPost"][0]
		category := formData["contentPost"]

	}
}
