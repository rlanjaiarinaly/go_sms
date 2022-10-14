package main

import (
	"SMS/modem"
	"fmt"
	"strings"
)

func main() {
	/* 	mod := modem.Modem{
	   		Numero:   0,
	   		Operator: "telma",
	   	}
	   	sms := fmt.Sprintln("Hello \nYour message contains some lines")
	   	fmt.Println(mod.SendSMS(sms, "0341542314")) */
	modems, _ := modem.GetAllModem()
	telma := modems["TELMA"]
	message := "This is a 'message' with quotes"
	fmt.Println(telma.SendSMS(message, "0341542314"))
}

func formatSMS(SMS string) string {
	return strings.Replace(SMS, "'", "\\'", -1)
}
