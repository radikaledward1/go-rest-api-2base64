package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/radikaledward1/golang-rest-api-2base64/routes"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)
	routes.RoutesRegister(router)
	fmt.Println("Servidor corriendo en http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
