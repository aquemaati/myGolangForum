package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/model"
)

func NewPost(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		index, err := UserAuthParser(r, db, w)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("ici aussie")
		cats, err := model.FetchCat(db)
		if err != nil {
			http.Error(w, "Could not get category infos: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		index.Cat = cats

		err = tpl.ExecuteTemplate(w, "createPost.html", index)
		if err != nil {
			http.Error(w, "Error rendering page: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})
}
