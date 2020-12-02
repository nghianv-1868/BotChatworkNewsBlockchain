package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"botnews/routes"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func main() {

	//Load env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		fmt.Println("Error loading .env file")
	}

	fmt.Println("Starting Server")
	if args := os.Args; len(args) > 1 && args[1] == "-register" {
		go routes.RegisterWebhook()
	}
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Welcome!\n")
	})
	router.POST("/webhook/chatwork", routes.HandleChatworkWebhook)
	//Listen to crc check and handle
	router.GET("/twitter/crc", routes.CrcCheck)
	//Listen to webhook event and handle
	router.GET("/twitter/webhook", routes.HandleTwitterWebhook)

	// c := cron.New()

	log.Fatal(http.ListenAndServe(":9090", router))
}
