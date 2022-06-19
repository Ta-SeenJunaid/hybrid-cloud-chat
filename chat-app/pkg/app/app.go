package app

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Ta-SeenJunaid/hybrid-cloud-chat/chat-app/pkg/apis"
	"github.com/gorilla/mux"
)

var tpl = template.Must(template.ParseFiles("pkg/app/index.html"))

var Messages []apis.Message

func initializeDatabase() {
	Messages = make([]apis.Message, 0)
}

func ReceiveMessages(w http.ResponseWriter, r *http.Request) {
	// /receive/ for print
	err := tpl.Execute(w, Messages)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("GET /")

}

func SendMessages(w http.ResponseWriter, r *http.Request) {
	message := r.FormValue("message")
	if len(message) > 0 {
		err := NatsConnection.Publish(Receiver, []byte(message))
		if err != nil {
			log.Fatalln(err)
		}
		Messages = append(Messages, apis.Message{
			Author: Sender,
			Body:   message,
			Time:   time.Now().Format("2022-06-02 02:23:45"),
		})
	}

	err := tpl.Execute(w, Messages)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("POST /")
}

func Run() {
	InitializeNats()
	initializeDatabase()

	r := mux.NewRouter()
	r.HandleFunc("/", ReceiveMessages).Methods(http.MethodGet)
	r.HandleFunc("/", SendMessages).Methods(http.MethodPost)

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "8000"
	}
	srv := &http.Server{
		Handler: r,
		Addr:    ":" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
