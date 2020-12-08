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

	"botnews/stream"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

type InfoUser struct {
	Data DataInfoUser `json:"data"`
}

type DataInfoUser struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
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
		messBody := strings.Split(mess, "add_tag=")[1]
		messBody = strings.ReplaceAll(messBody, " ", "")
		listTags := strings.Split(messBody, ",")
		messRes := `[info][title]Tags[/title]- Add success: ` + strings.Join(listTags, " , ") + `[/info]`
		sendMessage(messRes)
		addToFile(listTags, ".tags")

	} else if strings.Contains(mess, "remove_tag=") {
		messBody := strings.Split(mess, "remove_tag=")[1]
		messBody = strings.ReplaceAll(messBody, " ", "")
		listTags := strings.Split(messBody, ",")
		readData, err := ioutil.ReadFile(".tags")
		listTagsSuccess := []string{}
		listTagsNotContains := []string{}
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range listTags {
			if containsArrayString(strings.Split(string(readData), ","), v) {
				listTagsSuccess = append(listTagsSuccess, v)
			} else {
				listTagsNotContains = append(listTagsNotContains, v+"(not contains tags)")
			}
		}
		messRes := `[info][title]Tags[/title]- Remove success: ` + strings.Join(listTagsSuccess, " , ") + `
- Remove error: ` + strings.Join(listTagsNotContains, " , ") + `[/info]`
		sendMessage(messRes)
		removbeToFile(listTags, ".tags")
	} else if strings.Contains(mess, "add_follow=") {
		messBody := strings.Split(mess, "add_follow=")[1]
		messBody = strings.ReplaceAll(messBody, " ", "")
		listUsername := strings.Split(messBody, ",")
		listFollow := []string{}
		listFollowUsername := []string{}
		listNotFound := []string{}
		for _, v := range listUsername {
			if getUserIdFromUsername(v) != "" {
				listFollow = append(listFollow, getUserIdFromUsername(v))
				listFollowUsername = append(listFollowUsername, v)
			} else {
				listNotFound = append(listNotFound, v+"(not found)")
			}
		}

		messRes := `[info][title]Follow[/title]- Add success: ` + strings.Join(listFollowUsername, " , ") + `
- Add error: ` + strings.Join(listNotFound, " , ") + `[/info]`
		sendMessage(messRes)
		addToFile(listFollow, ".following")
	} else if strings.Contains(mess, "remove_follow=") {
		messBody := strings.Split(mess, "remove_follow=")[1]
		messBody = strings.ReplaceAll(messBody, " ", "")
		listRemove := strings.Split(messBody, ",")
		listUserIdRemove := []string{}
		listUsernameRemove := []string{}
		listUsernameNotFound := []string{}
		readData, err := ioutil.ReadFile(".following")
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range listRemove {
			if getUserIdFromUsername(v) != "" {
				if !containsArrayString(strings.Split(string(readData), ","), getUserIdFromUsername(v)) {
					listUsernameNotFound = append(listUsernameNotFound, v+"(not contains following)")
				} else {
					listUserIdRemove = append(listUserIdRemove, getUserIdFromUsername(v))
					listUsernameRemove = append(listUsernameRemove, v)
				}

			} else {
				listUsernameNotFound = append(listUsernameNotFound, v+"(not found)")
			}
		}
		messRes := `[info][title]Follow[/title]- Remove success: ` + strings.Join(listUsernameRemove, " , ") + `
- Remove error: ` + strings.Join(listUsernameNotFound, " , ") + `[/info]`
		sendMessage(messRes)
		removbeToFile(listUserIdRemove, ".following")
	}
}

func addToFile(_list []string, _nameFile string) {
	readData, err := ioutil.ReadFile(_nameFile)
	if err != nil {
		fmt.Println(err)
	}
	arrData := []string{}
	if string(readData) != "" {
		arrData = strings.Split(string(readData), ",")
	}
	for _, v := range _list {
		if !containsArrayString(arrData, v) {
			arrData = append(arrData, v)
		}
	}

	dataSave := []byte(strings.Join(arrData, ","))
	err1 := ioutil.WriteFile(_nameFile, dataSave, 0777)
	if err1 != nil {
		fmt.Println(err1)
	}
	restartStream()
}
func removbeToFile(_list []string, _nameFile string) {
	readData, err := ioutil.ReadFile(_nameFile)
	if err != nil {
		fmt.Println(err)
	}
	arrData := strings.Split(string(readData), ",")

	for _, v := range _list {
		if index := indexContainsArrayString(arrData, v); index != -1 {
			arrData = append(arrData[:index], arrData[index+1:]...)
		}
	}
	dataSave := []byte(strings.Join(arrData, ","))
	err1 := ioutil.WriteFile(_nameFile, dataSave, 0777)
	if err1 != nil {
		fmt.Println(err1)
	}
	restartStream()
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

func getUserIdFromUsername(_username string) string {
	// URL request
	reqURL := "https://api.twitter.com/2/users/by/username/" + _username

	token := os.Getenv("BEARER_TOKEN")

	client := &http.Client{}

	request, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var info InfoUser
	error := json.Unmarshal(body, &info)
	if error != nil {
		return ""
	}
	return info.Data.Id
}
func getUsernameFromUserId(_userId string) string {
	// URL request
	reqURL := "https://api.twitter.com/2/users/" + _userId

	token := os.Getenv("BEARER_TOKEN")

	client := &http.Client{}

	request, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	var info InfoUser
	error := json.Unmarshal(body, &info)
	if error != nil {
		return ""
	}
	fmt.Println(info)
	return info.Data.Username
}

func restartStream() {
	stream.StreamTwitter.Stop()
	stream.CreateStreamTwitter()
	go stream.Demux.HandleChan(stream.StreamTwitter.Messages)
}

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
