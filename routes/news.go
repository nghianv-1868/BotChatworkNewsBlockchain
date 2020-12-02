package routes

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	_ "reflect"

	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"

	"github.com/dghubble/oauth1"
)

type Post struct {
	Tag     string `json:"tag"`
	Link    string `json:"link"`
	Created string `json:"created"`
}

//Struct to parse webhook load
type WebhookLoad struct {
	UserId           string  `json:"for_user_id"`
	TweetCreateEvent []Tweet `json:"tweet_create_events"`
}

//Struct to parse tweet
type Tweet struct {
	Id    int64
	IdStr string `json:"id_str"`
	User  User
	Text  string
}

//Struct to parse user
type User struct {
	Id     int64
	IdStr  string `json:"id_str"`
	Name   string
	Handle string `json:"screen_name"`
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

func CrcCheck(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	//Set response header to json type
	writer.Header().Set("Content-Type", "application/json")
	//Get crc token in parameter
	token := request.URL.Query()["crc_token"]
	if len(token) < 1 {
		fmt.Fprintf(writer, "No crc_token given")
		return
	}

	//Encrypt and encode in base 64 then return
	h := hmac.New(sha256.New, []byte(os.Getenv("CONSUMER_SECRET")))
	h.Write([]byte(token[0]))
	encoded := base64.StdEncoding.EncodeToString(h.Sum(nil))
	//Generate response string map
	response := make(map[string]string)
	response["response_token"] = "sha256=" + encoded
	//Turn response map to json and send it to the writer
	responseJson, _ := json.Marshal(response)
	fmt.Fprintf(writer, string(responseJson))
}

func CreateClient() *http.Client {
	//Create oauth client with consumer keys and access token
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("ACCESS_TOKEN_KEY"), os.Getenv("ACCESS_TOKEN_SECRET"))

	return config.Client(oauth1.NoContext, token)

}

func RegisterWebhook() {
	fmt.Println("Registering webhook...")
	httpClient := CreateClient()

	//Set parameters
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/webhooks.json"
	values := url.Values{}
	values.Set("url", os.Getenv("APP_URL")+"/twitter/crc")
	fmt.Println(os.Getenv("APP_URL"))

	//Make Oauth Post with parameters
	resp, _ := httpClient.PostForm(path, values)
	defer resp.Body.Close()
	//Parse response and check response
	body, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		panic(err)
	}
	fmt.Println(data)
	fmt.Println("Webhook id of " + data["id"].(string) + " has been registered")
	SubscribeWebhook()
}

func SubscribeWebhook() {
	fmt.Println("Subscribing webapp...")
	client := CreateClient()
	path := "https://api.twitter.com/1.1/account_activity/all/" + os.Getenv("WEBHOOK_ENV") + "/subscriptions.json"
	resp, _ := client.PostForm(path, nil)
	fmt.Println(resp)
	body, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	//If response code is 204 it was successful
	if resp.StatusCode == 204 {
		fmt.Println("Subscribed successfully")
	} else if resp.StatusCode != 204 {
		fmt.Println("Could not subscribe the webhook. Response below:")
		fmt.Println(string(body))
	}
}

func HandleChatworkWebhook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}

func HandleTwitterWebhook(writer http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	fmt.Println("Handler called")
	//Read the body of the tweet
	body, _ := ioutil.ReadAll(request.Body)
	//Initialize a webhok load obhject for json decoding
	var load WebhookLoad
	err := json.Unmarshal(body, &load)
	if err != nil {
		fmt.Println("An error occured: " + err.Error())
	}
	//Check if it was a tweet_create_event and tweet was in the payload and it was not tweeted by the bot
	if len(load.TweetCreateEvent) < 1 || load.UserId == load.TweetCreateEvent[0].User.IdStr {
		return
	}
	//Send Hello world as a reply to the tweet, replies need to begin with the handles
	//of accounts they are replying to
	_, err = SendTweet("@"+load.TweetCreateEvent[0].User.Handle+" Hello World", load.TweetCreateEvent[0].IdStr)
	if err != nil {
		fmt.Println("An error occured:")
		fmt.Println(err.Error())
	} else {
		fmt.Println("Tweet sent successfully")
	}
}

func SendTweet(tweet string, reply_id string) (*Tweet, error) {
	fmt.Println("Sending tweet as reply to " + reply_id)
	//Initialize tweet object to store response in
	var responseTweet Tweet
	//Add params
	params := url.Values{}
	params.Set("status", tweet)
	params.Set("in_reply_to_status_id", reply_id)
	//Grab client and post
	client := CreateClient()
	resp, err := client.PostForm("https://api.twitter.com/1.1/statuses/update.json", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	//Decode response and send out
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	err = json.Unmarshal(body, &responseTweet)
	if err != nil {
		return nil, err
	}
	return &responseTweet, nil
}
