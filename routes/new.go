package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "reflect"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"

	"botNews/setupGoogleAPI"

	res "botNews/utils"
)

type New struct {
	tag     string `json:"tag"`
	link    string `json:"link"`
	created string `json:"created"`
}

func GetAllNews(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := setupGoogleAPI.GetClient(config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	spreadsheetID := "1q2df2VfVYnHYblqivWj5WUe5nxOh-Yj8RYxj6nUcpUU"
	readRange := "Class Data!A2:E"
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	news := make(map[string][]New)
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		res.JSON(w, 200, "No data found.")
	} else {
		fmt.Println("oke")
		var listNews []New
		for _, row := range resp.Values {
			tag, _ := row[0].(string)
			link, _ := row[4].(string)
			a := &New{tag: tag, link: link, created: time.Now().String()}
			e, err := json.MarshalIndent(a, "", "  ")
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(e))
			listNews = append(listNews, *a)
		}
		news["tag"] = listNews
	}
	fmt.Println(news)
	res.JSON(w, 200, news)
}
