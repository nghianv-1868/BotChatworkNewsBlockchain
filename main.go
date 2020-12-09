package main

import (
	"fmt"
	"log"
	"net/http"

	"botnews/routes"
	"botnews/stream"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func main() {

	//Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		fmt.Println("Error loading .env file")
	}

	// Create stream twitter
	stream.CreateStreamTwitter()
	go stream.Demux.HandleChan(stream.StreamTwitter.Messages)

	// Create server http
	fmt.Println("Starting Server")
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Welcome!\n")
	})
	router.POST("/chatwork/webhook", routes.HandleChatworkWebhook)
	log.Fatal(http.ListenAndServe(":6868", router))
}
