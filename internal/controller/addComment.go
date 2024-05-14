package controller

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/google/uuid"
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
		commId := uuid.New().String()

		//recuperer aussi id du post dans le html
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)
		postId := getFirstValue(formData, "postId")
		content := getFirstValue(formData, "content")

		fmt.Println("add comments paramewters")
		fmt.Println("user", index.UserID, "commId", commId, "postId", postId, "content", content)

		//redirect to index/
		err = SubmitComment(db, commId, postId, index.UserID, content)
		if err != nil {
			http.Error(w, "AddComment failed", http.StatusInternalServerError)
		}
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

func SubmitComment(db *sql.DB, commId, postId, userId, content string) error {
	if content == "" {
		return fmt.Errorf("can't post an empty comment")
	}

	stmt, err := db.Prepare("INSERT INTO Comments(id, postId, userId, content) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(commId, postId, userId, content)
	if err != nil {
		return err
	}

	return nil
}
