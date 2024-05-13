package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/aquemaati/myGolangForum.git/internal/model"
)

func FilterByReact(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("jeloooooooo")
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)
		userID, _ := r.Context().Value(middleware.UserIdContextKey).(string)
		fmt.Println(formData)
		react := formData["filter"][0]

		index := Index{}

		fmt.Println("Hello")
		posts, err := model.FetchPostsReactedByUser(db, userID, &react)
		if err != nil {
			http.Error(w, "could not get posts infos "+err.Error(), http.StatusInternalServerError)
			log.Panicln(err)
			return
		}

		fmt.Println(posts)
		index.Posts = posts
		cats, err := model.FetchCat(db)
		if err != nil {
			http.Error(w, "could not get cat infos "+err.Error(), http.StatusInternalServerError)
		}
		index.Cat = cats

		err = tpl.ExecuteTemplate(w, "index.html", index)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
