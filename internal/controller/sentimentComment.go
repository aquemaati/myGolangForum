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

func SentimentComment(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(middleware.UserIdContextKey).(string)
		if !ok {
			log.Println("could not find userid", userID)
		}
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)

		//postId
		commentId := formData["commentId"][0]
		//sentiment
		sentiment := formData["sentiment"][0]

		commentIdInt, err := strconv.Atoi(commentId)
		if err != nil {
			panic(err)
		}

		// Vérifier l'existence du sentiment actuel
		existingSentiment, err := model.GetUserSentimentComment(db, userID, commentIdInt)
		if err != nil && err != model.ErrSentimentNotFound {
			http.Error(w, "Failed to get user sentiment: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Si le sentiment actuel est le même que celui demandé, supprimer
		if existingSentiment == sentiment {
			err := model.RemoveSentomentComment(db, userID, commentIdInt)
			if err != nil {
				http.Error(w, "Could not remove sentiment: "+err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			// Ajouter ou mettre à jour le sentiment
			err = model.AddSentiimentComment(db, userID, commentIdInt, sentiment)
			if err != nil {
				http.Error(w, "Could not add/update sentiment: "+err.Error(), http.StatusInternalServerError)
				return
			}
		}

		// Rediriger après succès
		http.Redirect(w, r, "/", http.StatusFound)
	})
}
