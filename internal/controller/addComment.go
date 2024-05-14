package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/gofrs/uuid"
)

func AddComment(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and get user details
		index, err := UserAuthParser(r, db, w)
		if err != nil {
			log.Println("Authentication error:", err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
			return
		}

		// Generate a new UUID for the postId
		commId, err := uuid.NewV4()
		if err != nil {
			log.Println("Error generating UUID:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		//recuperer aussi id du post dans le html
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)
		postId := getFirstValue(formData, "postId")
		content := getFirstValue(formData, "content")

		fmt.Println("add comments paramewters")
		fmt.Println("user", index.UserID, "commId", commId, "postId", postId, "content", content)
		tpl.ExecuteTemplate(w, "index.html", nil)
	}
}

func SubmitComment(db *sql.DB, commId, postId, userId, content string) error {

	return nil
}
