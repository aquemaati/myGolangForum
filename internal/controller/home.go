package controller

import (
	"database/sql"
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
		userID, ok := r.Context().Value(middleware.UserIdContextKey).(string)
		if !ok {
			// Handle the case where the userID is missing or of an incorrect type
			log.Println("User ID not found in context or is of incorrect type")
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			return
		}

		index := Index{UserID: userID}

		// Get user public info from db
		user, err := model.GetUserPublicById(db, userID)
		if err != nil {
			http.Error(w, "Could not get user infos: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		index.UserInfos = user

		// Obtain all the posts info in the db
		posts, err := model.FetchExtendedPostsWithComments(db, nil, nil)
		if err != nil {
			http.Error(w, "Could not get posts infos: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		index.Posts = posts

		// Get categories from db
		cats, err := model.FetchCat(db)
		if err != nil {
			http.Error(w, "Could not get category infos: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		index.Cat = cats

		// Render the page
		err = tpl.ExecuteTemplate(w, "index.html", index)
		if err != nil {
			http.Error(w, "Error rendering page: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})
}
