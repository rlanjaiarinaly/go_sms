package main

import (
	"SMS/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	// modems, _ := modem.GetAllModem()
	router := mux.NewRouter()

	smsC := service.NewSMSController()
	router.HandleFunc("/send_sms", smsC.Send).Methods("POST")
	fmt.Println("Server listening on port 3000")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
