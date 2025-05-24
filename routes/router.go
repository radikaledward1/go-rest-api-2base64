package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/radikaledward1/golang-rest-api-2base64/services"
)

func RoutesRegister(router *mux.Router) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Welcome to File to Base64 API 🎉")
	})
	router.HandleFunc("/document", services.GetDocument).Methods("GET")
}
