package service

import (
	"SMS/models"
	"fmt"
	"net/http"
)

type SMSController struct {
	SMSService models.SMSService
}

func NewSMSController() *SMSController {
	return &SMSController{
		SMSService: models.NewSMSService(),
	}
}

func (smsC *SMSController) Send(w http.ResponseWriter, r *http.Request) {
	sms := models.SMS{}
	if err := Helper(&sms, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err := smsC.SMSService.Send(&sms)
	if err != nil {
		fmt.Fprintln(w, err)
	}
}
