package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/aquemaati/myGolangForum.git/internal/model"
)

func FilteredHome(db *sql.DB, tpl *template.Template) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)

		// Get category from form data safely
		category := getFirstValue(formData, "category")
		userId := getFirstValue(formData, "userId")

		var cat *string
		if category == "" {
			cat = nil // Directly assign nil to the pointer
		} else {
			cat = &category // Assign the address of the 'category' variable to the pointer
		}

		var uid *string
		if userId == "" {
			uid = nil // Directly assign nil to the pointer
		} else {
			uid = &userId // Assign the address of the 'category' variable to the pointer
		}

		posts, err := model.FetchExtendedPostsWithComments(db, uid, cat)
		if err != nil {
			http.Error(w, "Could not get posts info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		index := Index{
			Posts: posts,
		}

		fmt.Println(index.Posts)
		cats, err := model.FetchCat(db)
		if err != nil {
			http.Error(w, "Could not get category info: "+err.Error(), http.StatusInternalServerError)
			return
		}
		index.Cat = cats

		err = tpl.ExecuteTemplate(w, "index.html", index)
		if err != nil {
			http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
			return
		}
	})
}

// Utility function to safely get the first item from form data
func getFirstValue(data map[string][]string, key string) string {
	values, exists := data[key]
	if exists && len(values) > 0 {
		return values[0]
	}
	return ""
}
