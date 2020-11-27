package routes

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	_ "reflect"

	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type Post struct {
	Tag     string `json:"tag"`
	Link    string `json:"link"`
	Created string `json:"created"`
}

func GetAllNews(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func sendMessage(message string) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// URL request
	reqURL := "https://api.chatwork.com/v2/rooms/206069293/messages?body=" + url.QueryEscape(message)

	token := os.Getenv("TOKEN_CHATWORK_BOT")

	client := &http.Client{}

	request, err := http.NewRequest("POST", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("X-ChatworkToken", token)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func PostNews(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sendMessage("Hello Everyone")
}
