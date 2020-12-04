package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	_ "reflect"
	"strings"

	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type Post struct {
	Tag     string `json:"tag"`
	Link    string `json:"link"`
	Created string `json:"created"`
}

//Struct of WebhookChatwork
type WebhookChatwork struct {
	Webhook_setting_id string        `json:"webhook_setting_id"`
	Webhook_event_type string        `json:"webhook_event_type"`
	Webhook_event_time int64         `json:"webhook_event_time"`
	Webhook_event      Webhook_event `json:"webhook_event"`
}

//Struct of Webhook_event
type Webhook_event struct {
	From_account_id int64  `json:"From_account_id"`
	To_account_id   int64  `json:"to_account_id"`
	Room_id         int64  `json:"room_id`
	Message_id      string `json:"message_id"`
	Body            string `json:"body"`
	Send_time       int64  `json:"send_time"`
	Update_time     int64  `json:"update_time"`
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

func HandleChatworkWebhook(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// ValidateRequestChatwork(r)
	body, _ := ioutil.ReadAll(r.Body)
	//Initialize a webhok load obhject for json decoding
	var load WebhookChatwork
	err := json.Unmarshal(body, &load)
	if err != nil {
		fmt.Println("An error occured: " + err.Error())
	}
	mess := string(load.Webhook_event.Body)
	if strings.Contains(mess, "add_tag=") {
		messTags := strings.Split(mess, "add_tag=")[1]
		messTags = strings.ReplaceAll(messTags, " ", "")
		tags := strings.Split(messTags, ",")
		readData, err := ioutil.ReadFile(".following")
		if err != nil {
			fmt.Println(err)
		}
		arrData := strings.Split(string(readData), ",")

		for _, v := range tags {
			if !containsArrayString(arrData, v) {
				arrData = append(arrData, v)
			}
		}
		dataSave := []byte(strings.Join(arrData, ","))
		err1 := ioutil.WriteFile(".following", dataSave, 0777)
		if err1 != nil {
			fmt.Println(err1)
		}
	} else if strings.Contains(mess, "remove_tag=") {
		messTags := strings.Split(mess, "remove_tag=")[1]
		messTags = strings.ReplaceAll(messTags, " ", "")
		tags := strings.Split(messTags, ",")
		readData, err := ioutil.ReadFile(".following")
		if err != nil {
			fmt.Println(err)
		}
		arrData := strings.Split(string(readData), ",")

		for _, v := range tags {
			if index := indexContainsArrayString(arrData, v); index != -1 {
				arrData = append(arrData[:index], arrData[index+1:]...)
			}
		}
		dataSave := []byte(strings.Join(arrData, ","))
		err1 := ioutil.WriteFile(".following", dataSave, 0777)
		if err1 != nil {
			fmt.Println(err1)
		}
	}
}

func containsArrayString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func indexContainsArrayString(s []string, str string) int {
	for i, v := range s {
		if v == str {
			return i
		}
	}

	return -1
}

// func RemoveDuplicates(s []string) []string {
// 	m := make(map[string]bool)
// 	for _, item := range s {
// 		if _, ok := m[item]; ok {
// 			// duplicate item
// 			fmt.Println(item, "is a duplicate")
// 		} else {
// 			m[item] = true
// 		}
// 	}

// 	var result []string
// 	for item, _ := range m {
// 		result = append(result, item)
// 	}
// 	return result
// }

// func ValidateRequestChatwork(r *http.Request) {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 		fmt.Println("Error loading .env file")
// 	}
// 	// get chatwork_signature
// 	chatwork_signature := string(r.Header.Get("X-Chatworkwebhooksignature"))

// 	// hash body request
// 	decodeToken, _ := base64.StdEncoding.DecodeString(os.Getenv("TOKEN_WEBHOOK_CHATWORK"))
// 	body, err := ioutil.ReadAll(r.Body)
// 	h := hmac.New(sha256.New, []byte(decodeToken))
// 	h.Write(body)
// 	sha := hex.EncodeToString(h.Sum(nil))

// 	// compare values
// 	fmt.Println(string(chatwork_signature) == sha)
// }
