package controller

import (
	"fmt"
	"net/http"
)

// HomeHandler handles the root path
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World")
}
