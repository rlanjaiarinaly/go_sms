package models

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
)

type SMS struct {
	ID       int
	Number   string `schema:"numero"`
	Content  string `schema:"contenu"`
	Operator string
}
type smsValidator struct {
	SMSService
}
type smsService struct {
	SMSService
}

type SMSService interface {
	Send(sms *SMS) error
}

var (
	ErrNumberValue  = errors.New("number : the number is invalid")
	ErrEmptyContent = errors.New("sms : the content is empty")
	operators       = map[rune]string{'2': "ORANGE", '4': "TELMA", '8': "TELMA", '3': "AIRTEL"}
)

func NewSMSService() SMSService {
	return &smsService{
		SMSService: &smsValidator{},
	}
}

func (smsV *smsValidator) setOperatorName(sms *SMS) error {
	op := string(sms.Number[1]) + string(sms.Number[2])
	switch op {
	case "32":
	case "33":
	case "34", "38":
	default:
		return ErrNumberValue
	}
	sms.Operator = operators[rune(op[1])]
	return nil
}

func (smsV *smsValidator) Send(sms *SMS) error {
	if sms, err := smsV.newSMS(sms); err != nil {
		return err
	} else {
		return sms.send()
	}
}
func (smsV *smsValidator) validateNumber(sms *SMS) error {
	sms.Number = strings.Replace(strings.Replace(sms.Number, " ", "", -1), "+261", "0", 1)
	if len(sms.Number) != 10 {
		return ErrNumberValue
	}
	op := string(sms.Number[1]) + string(sms.Number[2])
	switch op {
	case "32":
	case "33":
	case "34", "38":
	default:
		return ErrNumberValue
	}
	return nil
}

func (smsV *smsValidator) normalizeContent(sms *SMS) error {
	if sms.Content == "" {
		return ErrEmptyContent
	}
	sms.Content = strings.Replace(sms.Content, "\"", "''", -1)
	return nil
}

// Get the SMS number from the output of the sms creating command
func getSMSID(out *bytes.Buffer) (int, error) {
	s := strings.TrimSpace(out.String())
	s = strings.Replace(s, " ", "", -1)
	splitted := strings.Split(s, "/")
	number, err := strconv.Atoi(splitted[len(splitted)-1])
	if err != nil {
		return -1, ErrFetchingSMS
	}
	return number, nil
}

type smsValFunc func(sms *SMS) error

func runSMSValFunc(sms *SMS, smsValFunc ...smsValFunc) error {
	for _, v := range smsValFunc {
		if err := v(sms); err != nil {
			return err
		}
	}
	return nil
}

func (smsV *smsValidator) newSMS(sms *SMS) (*SMS, error) {
	modems, err := GetAllModem()
	if err != nil {
		return nil, err
	}

	if err = runSMSValFunc(sms, smsV.validateNumber, smsV.setOperatorName, smsV.normalizeContent); err != nil {
		return nil, err
	}
	modem := modems[sms.Operator]
	if modem == nil {
		return nil, errors.New("modem : il n'y a pas de modem associé à l'operateur du numero branché")
	}
	if err = modem.NewSMS(sms); err != nil {
		return nil, err
	}
	return sms, nil
}

func (sms *SMS) send() error {
	cmd := exec.Command("mmcli", "--sms", strconv.Itoa(sms.ID), "--send")
	return cmd.Run()
}
