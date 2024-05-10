package controller

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/model"
)

func Test(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cat, err := model.FetchUniquePost(db, 2)
		if err != nil {
			http.Error(w, "could not get cat infos "+err.Error(), http.StatusInternalServerError)
		}

		err = tpl.ExecuteTemplate(w, "test.html", cat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
