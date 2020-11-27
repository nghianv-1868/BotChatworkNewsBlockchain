package main

import (
	"fmt"
	"log"
	"net/http"

	"botnews/routes"

	"github.com/julienschmidt/httprouter"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {

	router := httprouter.New()
	router.GET("/wirteNews", routes.PostNews)
	log.Fatal(http.ListenAndServe(":6379", router))
}
