package server

// import (
// 	"net/http"
// 	"testing"
// 	"time"
// )

// // Test server initialization
// func TestInitializeServer(t *testing.T) {
// 	// Arrange
// 	server := InitializeServer()

// 	// Act & Assert
// 	if server.Addr != ":9090" {
// 		t.Errorf("Expected server to listen on :8080 but got %v", server.Addr)
// 	}
// 	if server.ReadTimeout != 10*time.Second {
// 		t.Errorf("Expected ReadTimeout to be 10s but got %v", server.ReadTimeout)
// 	}
// 	if server.MaxHeaderBytes != 1<<20 {
// 		t.Errorf("Expected MaxHeaderBytes to be 1 MB but got %v", server.MaxHeaderBytes)
// 	}
// 	if _, ok := server.Handler.(http.Handler); !ok {
// 		t.Error("Expected handler to implement http.Handler")
// 	}
// }
