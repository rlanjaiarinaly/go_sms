package main

import (
	"SMS/modem"
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
	message := "This is a \"message\" , \n with quotes"
	messages := map[string]string{
		"0341542314": message,
	}
	for i, v := range messages {
		telma.SendSMSToRecipient(v, i)
	}
	// fmt.Println(telma.SendSMS(message, "0341542314"))

}
