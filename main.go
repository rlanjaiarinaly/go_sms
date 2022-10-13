package main

import (
	"SMS/modem"
	"fmt"
)

func main() {
	/* 	mod := modem.Modem{
	   		Numero:   0,
	   		Operator: "telma",
	   	}
	   	sms := fmt.Sprintln("Hello \nYour message contains some lines")
	   	fmt.Println(mod.SendSMS(sms, "0341542314")) */
	fmt.Println(modem.GetAllModem())
}
