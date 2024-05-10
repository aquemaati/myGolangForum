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

type Index struct {
	Cat       []model.Categorie
	Posts     []model.PostInfo
	UserID    string
	UserInfos model.UserPublic
}

// HomeHandler handles the root path
func Home(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Implémentation de la page d'accueil
		// Exemple simple d'une réponse HTML

		userID, ok := r.Context().Value(middleware.UserIdContextKey).(string)
		index := Index{}
		if !ok {
			// Handle the case where the userID is missing or of an incorrect type
			log.Println("User ID not found in context", http.StatusUnauthorized)
		} else {

			user, err := model.GetUserPublicById(db, userID)
			if err != nil {
				http.Error(w, "could not get user infos "+err.Error(), http.StatusInternalServerError)
				log.Panicln(err)
				return
			}
			index.UserInfos = user
		}

		fmt.Println("Retrieved User ID:", userID)
		index.UserID = userID

		posts, err := model.FetchExtendedPostsWithComments(db, nil, nil)
		if err != nil {
			http.Error(w, "could not get posts infos "+err.Error(), http.StatusInternalServerError)
			log.Panicln(err)
			return
		}
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
