package main

import (
	"fmt"
	"log"
	"net/http"

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

	// a := []string{"Defi", "A", "B"}
	// mydata := []byte(strings.Join(a, ","))
	// err1 := ioutil.WriteFile(".following", mydata, 0777)
	// if err1 != nil {
	// 	fmt.Println(err1)
	// }

	// data, err := ioutil.ReadFile(".following")
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Print(strings.Split(string(data), "-"))
	// fmt.Print(strings.Join(strings.Split(string(data), "-"), "-"))

	fmt.Println("Starting Server")
	router := httprouter.New()
	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Welcome!\n")
	})
	router.POST("/chatwork/webhook", routes.HandleChatworkWebhook)

	log.Fatal(http.ListenAndServe(":9090", router))
}
