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

func UniquePost(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)

		postId := formData["postId"][0]
		fmt.Println(postId)

		index := Index{}

		posts, err := model.FetchUniquePost(db, postId)
		if err != nil {
			http.Error(w, "could not get posts infos "+err.Error(), http.StatusInternalServerError)
			log.Panicln(err)
			return
		}
		index.Posts = append(index.Posts, posts)

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
