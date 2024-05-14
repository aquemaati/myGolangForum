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
	Auth      bool
	Filtered  bool
}

// HomeHandler handles the root path
func Home(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		index, err := UserAuthParser(r, db, w)
		if err != nil {
			log.Println(err)
		}

		// Fetch all posts info from the database, which is visible to both authenticated users and public visitors
		posts, err := model.FetchExtendedPostsWithComments(db, nil, nil)
		if err != nil {
			http.Error(w, "Could not get posts infos: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		index.Posts = posts

		// Fetch categories from the database
		cats, err := model.FetchCat(db)
		if err != nil {
			http.Error(w, "Could not get category infos: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		index.Cat = cats

		// Render the page with data, tailored to the authentication state of the user
		err = tpl.ExecuteTemplate(w, "index.html", index)
		if err != nil {
			http.Error(w, "Error rendering page: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
	})
}

func UserAuthParser(r *http.Request, db *sql.DB, w http.ResponseWriter) (Index, error) {
	index := Index{}
	userID, ok := r.Context().Value(middleware.UserIdContextKey).(string)
	if ok {

		index.UserID = userID
		index.Auth = true
		user, err := model.GetUserPublicById(db, userID)
		if err != nil {
			http.Error(w, "Could not get user infos: "+err.Error(), http.StatusInternalServerError)
			log.Println(err)
		}
		index.UserInfos = user
	} else {
		log.Println("User ID not found in context or is of incorrect type; proceeding as public visitor")
		index.Auth = false
	}
	return index, nil
}
