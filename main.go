package main

import (
	"fmt"
	"log"
	"net/http"

	"botNews/routes"

	"github.com/julienschmidt/httprouter"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {

	router := httprouter.New()
	router.GET("/", routes.GetAllNews)
	log.Fatal(http.ListenAndServe(":6379", router))
}
