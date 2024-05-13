package controller

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"

	"github.com/aquemaati/myGolangForum.git/internal/middleware"
	"github.com/gofrs/uuid"
)

func SubmitPostHandler(db *sql.DB, tpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Authenticate and get user details
		index, err := UserAuthParser(r, db, w)
		if err != nil {
			log.Println("Authentication error:", err)
			http.Error(w, "Authentication failed", http.StatusInternalServerError)
			return
		}

		// Generate a new UUID for the postId
		postId, err := uuid.NewV4()
		if err != nil {
			log.Println("Error generating UUID:", err)
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		// Retrieve form data
		formData := r.Context().Value(middleware.FormDataKey).(map[string][]string)
		titlePost := formData["titlePost"][0]
		descriptionPost := formData["descriptionPost"][0]
		categories := formData["category"] // Assuming checkboxes are used and name is 'category'

		// Insert the post and associate categories
		if err := SubmitPost(db, postId.String(), index.UserID, titlePost, descriptionPost, categories); err != nil {
			log.Println("Error submitting post:", err)
			http.Error(w, "Error submitting post", http.StatusInternalServerError)
			return
		}

		// Redirect or render success message/template
		http.Redirect(w, r, "/post-success", http.StatusSeeOther)
	}
}

func SubmitPost(db *sql.DB, postId, userId, title, description string, categoryNames []string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Insert the post into the Posts table
	_, err = tx.Exec("INSERT INTO Posts (id, userId, title, description) VALUES (?, ?, ?, ?)", postId, userId, title, description)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Prepare the query to fetch category IDs based on their names
	stmt, err := tx.Prepare("SELECT id FROM Categories WHERE name = ?")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	// For each category name, find its ID and insert into PostCategories
	for _, categoryName := range categoryNames {
		var categoryId int
		if err := stmt.QueryRow(categoryName).Scan(&categoryId); err != nil {
			tx.Rollback()
			return err
		}

		if _, err := tx.Exec("INSERT INTO PostCategories (postId, categoryId) VALUES (?, ?)", postId, categoryId); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
