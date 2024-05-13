package controller

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/aquemaati/myGolangForum.git/internal/model"
)

func SentimentPost(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIdContextKey).(string)
		if !ok {
			log.Println("could not find userid", userID)
		}
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)

		//postId
		postId := formData["postId"][0]
		//sentiment
		sentiment := formData["sentiment"][0]

		postIdInt, err := strconv.Atoi(postId)
		if err != nil {
			panic(err)
		}

		// Vérifier l'existence du sentiment actuel
		existingSentiment, err := model.GetUserSentiment(db, userID, postIdInt)
		if err != nil && err != model.ErrSentimentNotFound {
			http.Error(w, "Failed to get user sentiment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Si le sentiment actuel est le même que celui demandé, supprimer
		if existingSentiment == sentiment {
			err := model.RemovePostLike(db, userID, postIdInt)
			if err != nil {
				http.Error(w, "Could not remove sentiment: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// Ajouter ou mettre à jour le sentiment
			err = model.AddPostLike(db, userID, postIdInt, sentiment)
			if err != nil {
				http.Error(w, "Could not add/update sentiment: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Rediriger après succès
		referrer := r.Header.Get("Referer")
		if referrer == "" {
			referrer = "/" // Redirigez vers la page d'accueil si le referer est absent
		}
		http.Redirect(w, r, referrer, http.StatusFound)
	})
}
