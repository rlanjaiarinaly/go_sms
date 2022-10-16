package models

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// This is the class Modem containing the necessary informations inside the modem
type Modem struct {
	Numero   int
	Operator string
}

var (
	ErrCreatingSMS = errors.New("sms: there was an error while creating the sms")
	ErrSendingSMS  = errors.New("sms: there was an error while sending the sms")
	ErrNoOperator  = errors.New("sim: the operator can't be found")
	ErrFetchingSMS = errors.New("sms: failed fetching the sms")
)

// Fetch the list of the modems connected to the computer, note that mmcli doesn't support 4G and later dongle.
// This function returns a map of modem, the key is the operator of the sim inserted in the modem in uppercase.
func GetAllModem() (map[string]*Modem, error) {
	cmd := exec.Command("mmcli", "-L")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, err
	}
	modems := strings.Split(out.String(), "\n")
	modems = modems[:len(modems)-1]
	var op = make(map[string]*Modem)
	for _, v := range modems {
		infos := strings.Split(strings.TrimSpace(v), " ")
		paths := strings.Split(infos[0], "/")
		number, _ := strconv.Atoi(paths[len(paths)-1])
		cmd := exec.Command("mmcli", "-m", strconv.Itoa(number))
		var stdout bytes.Buffer
		cmd.Stdout = &stdout
		if err := cmd.Run(); err != nil {
			continue
		}
		operator, err := getOperatorName(&stdout)
		if err != nil {
			continue
		}
		op[strings.ToUpper(operator)] = &Modem{Numero: number, Operator: operator}
	}
	return op, nil
}

func getOperatorName(out *bytes.Buffer) (string, error) {
	regex := regexp.MustCompile("operator name: (?P<deb>[a-zA-Z]+)")
	match := regex.FindStringSubmatch(out.String())
	if len(match) < 2 {
		return "", ErrNoOperator
	}
	return match[1], nil
}

func (mod *Modem) NewSMS(sms *SMS) error {
	initiateCommand := fmt.Sprintf("--messaging-create-sms=number=%s,text=\"%s\"", sms.Number, sms.Content)
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("mmcli", "-m", strconv.Itoa(mod.Numero), initiateCommand)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return ErrCreatingSMS
	}
	ID, err := getSMSID(&stdout)
	if err != nil {
		return err
	}
	sms.ID = ID
	return nil
}
