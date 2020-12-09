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

type ResponseRoom struct {
	Room_id          int    `json:"room_id"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	Role             string `json:"role"`
	Sticky           bool   `json:"sticky"`
	Unread_num       int    `json:"unread_num"`
	Mention_num      int    `json:"mention_num"`
	Mytask_num       int    `json:"mytask_num"`
	Message_num      int    `json:"message_num"`
	File_num         int    `json:"file_num"`
	Task_num         int    `json:"task_num"`
	Icon_path        string `json:"icon_path"`
	Last_update_time int    `json:"last_update_time"`
	Description      string `json:"description"`
}

func sendMessage(message string) {
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
	updateWhenChangeFile()
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
	updateWhenChangeFile()
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

func updateWhenChangeFile() {
	stream.StreamTwitter.Stop()
	stream.CreateStreamTwitter()
	go stream.Demux.HandleChan(stream.StreamTwitter.Messages)

	dataTags, err := ioutil.ReadFile(".tags")
	if err != nil {
		fmt.Println(err)
	}

	dataFollowing, err := ioutil.ReadFile(".following")
	if err != nil {
		fmt.Println(err)
	}

	following := strings.Split(string(dataFollowing), ",")

	for i, v := range following {
		following[i] = getUsernameFromUserId(v)
	}

	// Update list following
	description := getDescriptionChatwork()
	fromF := strings.Index(description, "[info][title]List Following[/title] [") + 37
	toF := strings.Index(description, "] [/info][info][title]List Tags[/title]")
	description = description[:fromF] + strings.Join(following, ",") + description[(toF-2):]

	// Update list tags
	fromT := strings.Index(description, "[info][title]List Tags[/title] [") + 32
	description = description[:fromT] + string(dataTags) + "] [/info]"

	updateDescriptionChatwork(description)
}

func getDescriptionChatwork() string {
	// URL request
	reqURL := "https://api.chatwork.com/v2/rooms/206069293"

	token := os.Getenv("TOKEN_CHATWORK_BOT")

	client := &http.Client{}

	request, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("X-ChatworkToken", token)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	fmt.Println("Get Description:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	var res ResponseRoom
	error := json.Unmarshal(body, &res)
	if error != nil {
		return ""
	}
	return res.Description
}

func updateDescriptionChatwork(message string) {
	// URL request
	reqURL := "https://api.chatwork.com/v2/rooms/206069293?description=" + url.QueryEscape(message)

	token := os.Getenv("TOKEN_CHATWORK_BOT")

	client := &http.Client{}

	request, err := http.NewRequest("PUT", reqURL, nil)
	if err != nil {
		log.Fatalln(err)
	}

	request.Header.Set("X-ChatworkToken", token)

	resp, err := client.Do(request)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	fmt.Println("Update Description:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
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
