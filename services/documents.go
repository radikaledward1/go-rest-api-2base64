package services

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

func GetDocument(w http.ResponseWriter, r *http.Request) {
	// Log to console
	fmt.Println("Getting the Document!!!")

	// Create JSON response
	response := Response{
		Message: "Document service endpoint",
	}

	// Set response headers and send JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
