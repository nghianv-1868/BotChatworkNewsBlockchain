package stream

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/joho/godotenv"
)

var StreamTwitter *twitter.Stream
var Demux twitter.SwitchDemux

func CreateStreamTwitter() {

	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", os.Getenv("CONSUMER_KEY"), "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", os.Getenv("CONSUMER_SECRET"), "Twitter Consumer Secret")
	accessToken := flags.String("access-token", os.Getenv("ACCESS_TOKEN_KEY"), "Twitter Access Token")
	accessSecret := flags.String("access-secret", os.Getenv("ACCESS_TOKEN_SECRET"), "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter Client
	client := twitter.NewClient(httpClient)

	// Convenience Demux demultiplexed stream messages
	Demux = twitter.NewSwitchDemux()
	Demux.Tweet = func(tweet *twitter.Tweet) {
		if tweet.InReplyToStatusID == 0 {
			readData, err := ioutil.ReadFile(".following")
			arrData := strings.Split(string(readData), ",")
			if err != nil {
				fmt.Println(err)
			}
			if tweet.RetweetedStatus != nil && containsArrayString(arrData, tweet.User.IDStr) {
				messRes := "[info][title]https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr + " ( @" + tweet.User.ScreenName + " - Retweet ) [/title] " + tweet.RetweetedStatus.Text + "[/info]"
				sendMessage(messRes)
			} else if tweet.RetweetedStatus == nil && tweet.QuotedStatus != nil && containsArrayString(arrData, tweet.User.IDStr) {
				messRes := "[info][title]https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr + " ( @" + tweet.User.ScreenName + " - Quote ) [/title] " + tweet.Text + "[info]" + tweet.QuotedStatus.Text + "[/info][/info]"
				sendMessage(messRes)
			} else if tweet.RetweetedStatus == nil && tweet.QuotedStatus == nil && containsArrayString(arrData, tweet.User.IDStr) {
				messRes := "[info][title]https://twitter.com/" + tweet.User.ScreenName + "/status/" + tweet.IDStr + " ( @" + tweet.User.ScreenName + " - Tweet ) [/title] " + tweet.Text + "[/info]"
				sendMessage(messRes)
			}

		}

	}
	Demux.FriendsList = func(friendsList *twitter.FriendsList) {
		fmt.Println(friendsList)
	}
	Demux.DM = func(dm *twitter.DirectMessage) {
		fmt.Println(dm.SenderID)
	}
	Demux.Event = func(event *twitter.Event) {
		// fmt.Printf("%#v\n", event)
		fmt.Println("event", event.Event)
	}

	fmt.Println("Starting Stream...")

	listFollowing, err := ioutil.ReadFile(".following")
	fmt.Println("listFollowing: ", strings.Split(string(listFollowing), ","))
	if err != nil {
		fmt.Println(err)
	}

	listTags, err := ioutil.ReadFile(".tags")
	fmt.Println("listTags: ", strings.Split(string(listTags), ","))
	if err != nil {
		fmt.Println(err)
	}

	// FILTER
	filterParams := &twitter.StreamFilterParams{
		Follow:        strings.Split(string(listFollowing), ","),
		Track:         strings.Split(string(listTags), ","),
		StallWarnings: twitter.Bool(true),
	}
	StreamTwitter, err = client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}

	// go demux.HandleChan(stream.Messages)

	// // Wait for SIGINT and SIGTERM (HIT CTRL-C)
	// ch := make(chan os.Signal)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// log.Println(<-ch)

	// fmt.Println("Stopping Stream...")
	// stream.Stop()
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

	fmt.Println("SendMessage:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func containsArrayString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
