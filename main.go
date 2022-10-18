package main

import (
	"SMS/models"
	"SMS/service"
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	erreur := errors.New("initializing all the modems : ")
	for erreur != nil {
		fmt.Println(erreur)
		modems, erreur := models.GetAllModem()
		if erreur != nil {
			continue
		}
		for _, v := range modems {
			cmd := exec.Command("mmcli", "-m", strconv.Itoa(v.Numero), "-e")
			var stdout bytes.Buffer
			cmd.Stdout = &stdout
			erreur = cmd.Run()
			if erreur != nil {
				break
			}
			out := fmt.Sprintf("...%s %d Operator : %s", strings.Trim(stdout.String(), "\n"), v.Numero, v.Operator)
			fmt.Println(out)
		}
		if erreur != nil {
			continue
		}
		break
	}

	smsC := service.NewSMSController()
	router.HandleFunc("/send_sms", smsC.Send).Methods("POST")
	fmt.Println("Server listening on port 3000")
	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}
}
